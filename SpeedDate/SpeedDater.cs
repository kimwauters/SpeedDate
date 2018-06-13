﻿using System;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Text.RegularExpressions;
using SpeedDate.Configuration;
using SpeedDate.Interfaces;
using SpeedDate.Logging;
using SpeedDate.Network;
using SpeedDate.Network.Interfaces;
using SpeedDate.Plugin;
using SpeedDate.Plugin.Interfaces;

namespace SpeedDate
{
    public sealed class SpeedDater
    {
        private readonly string _configFile;
        private TinyIoCContainer _kernel;

        public event Action Started;
        public event Action Stopped;

        public bool IsStarted { get; set; } = false;

        public IPluginProvider PluginProver
        {
            get;
            private set;
        }

        public SpeedDater(string configFile)
        {
            _configFile = configFile;
        }

        public void Start()
        {
            SpeedDateConfig.FromXml(_configFile);
            var logger = LogManager.GetLogger("SpeedDate");

            _kernel = CreateKernel();

            var startable = _kernel.Resolve<ISpeedDateStartable>();
            startable.Started += () =>
            {
                IsStarted = true;
                Started?.Invoke();
            };
            startable.Stopped += () =>
            {
                IsStarted = false;
                Stopped?.Invoke();
            };

            PluginProver = _kernel.Resolve<IPluginProvider>();

            foreach (var plugin in _kernel.ResolveAll<IPlugin>())
            {
                if (SpeedDateConfig.Plugins.LoadAll || SpeedDateConfig.Plugins.PluginsNamespaces.Split(';').Any(ns => Regex.IsMatch(plugin.GetType().Namespace, WildCardToRegular(ns))))
                {
                    _kernel.BuildUp(plugin);
                    PluginProver.RegisterPlugin(plugin);
                }
            }

            foreach (var plugin in PluginProver.GetAll())
            {
                plugin.Loaded(PluginProver);
                logger.Info($"Loaded {plugin.GetType().Name}");
            }

            startable.Start();
        }

        public void Stop()
        {
            _kernel.TryResolve<ISpeedDateStartable>(out var startable);
            startable.Stop();
        }

        private static TinyIoCContainer CreateKernel()
        {
            try
            {
                //Register possible plugin-dependencies
                TinyIoCContainer.Current.Register<IClientSocket, ClientSocket>();
                TinyIoCContainer.Current.Register<IServerSocket, ServerSocket>();
                TinyIoCContainer.Current.Register<IPluginProvider, PluginProvider>();
                TinyIoCContainer.Current.Register<ILogger>((container, overloads, requestType) => LogManager.GetLogger(requestType.Name));

                //Register plugins
                foreach (var dllFile in
                    Directory.GetFiles(Path.GetDirectoryName(Assembly.GetExecutingAssembly().Location) ?? throw new InvalidOperationException(), "*.dll"))
                {
                    var assembly = Assembly.LoadFrom(dllFile);

                    foreach (var startableType in assembly.DefinedTypes.Where(info =>
                        !info.IsAbstract && !info.IsInterface && typeof(ISpeedDateStartable).IsAssignableFrom(info)))
                    {
                        var startableInstance = (ISpeedDateStartable)Activator.CreateInstance(startableType);

                        if(startableInstance is IServer)
                        {
                            TinyIoCContainer.Current.Register((container, overloads, requesttype) =>
                                (IServer) startableInstance);
                        }

                        if (startableInstance is IClient)
                        {
                            TinyIoCContainer.Current.Register((container, overloads, requesttype) =>
                                (IClient) startableInstance);
                        }
                        
                        TinyIoCContainer.Current.BuildUp(startableInstance);
                        TinyIoCContainer.Current.Register(startableInstance);
                    }

                    foreach (var pluginType in assembly.DefinedTypes.Where(info =>
                        !info.IsAbstract && !info.IsInterface && typeof(IPlugin).IsAssignableFrom(info)))
                    {
                        var pluginInstance = (IPlugin)Activator.CreateInstance(pluginType);
                        TinyIoCContainer.Current.Register(pluginInstance, pluginType.FullName);
                    }
                }
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex);
                throw;
            }

            return TinyIoCContainer.Current;
        }

        private static string WildCardToRegular(string value)
        {
            return "^" + Regex.Escape(value).Replace("\\*", ".*") + "$";
        }
    }
}
