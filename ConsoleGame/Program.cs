﻿using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using SpeedDate;
using SpeedDate.Client.Console.Example;

namespace ConsoleGame
{
    class Program
    {
        static void Main(string[] args)
        {
            Console.WriteLine($"Starting game with arguments: {string.Join(", ", args)}...");

            var server = new GameServer("gameserver.json");
            server.Start();

            server.ConnectedToMaster += () =>
            {
                Console.WriteLine("Connected to Master");
                server.Rooms.RegisterSpawnedProcess(
                    CommandLineArgs.SpawnId, 
                    CommandLineArgs.SpawnCode,
                    (controller, error) =>
                    {
                        Console.WriteLine(error ?? "Registered to Master");
                    });
            };

            Console.ReadLine();
        }
    }
}