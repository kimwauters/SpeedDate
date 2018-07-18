﻿using System.Collections.Generic;
using SpeedDate.Network;
using SpeedDate.Network.Utils.IO;

namespace SpeedDate.Packets.Lobbies
{
    /// <summary>
    /// This package represents current state of the lobby
    /// </summary>
    public class LobbyDataPacket : SerializablePacket
    {
        public LobbyState LobbyState;
        public string StatusText = "";
        public string GameMaster = "";
        public string CurrentUserUsername = "";

        public int LobbyId;
        public string LobbyName;
        public Dictionary<string, string> LobbyProperties;

        public int MaxPlayers;

        public Dictionary<string, LobbyMemberData> Players;
        public Dictionary<string, LobbyTeamData> Teams;
        public Dictionary<string, LobbyPropertyData> Controls;

        public byte[] AdditionalData;

        // Settings
        public bool EnableManualStart;
        public bool EnableReadySystem;
        public bool EnableTeamSwitching;


        public LobbyDataPacket()
        {
            // Just to avoid handling "null" cases
            Players = new Dictionary<string, LobbyMemberData>();
            Teams = new Dictionary<string, LobbyTeamData>();
            Controls = new Dictionary<string, LobbyPropertyData>();
        }

        public override void ToBinaryWriter(EndianBinaryWriter writer)
        {
            writer.Write((int)LobbyState);
            writer.Write(StatusText);
            writer.Write(GameMaster);
            writer.Write(CurrentUserUsername);

            writer.Write(LobbyId);
            writer.Write(LobbyName);
            writer.WriteDictionary(LobbyProperties);
            writer.Write(MaxPlayers);

            // Write additional data
            writer.Write(AdditionalData?.Length ?? 0);
            if (AdditionalData != null) writer.Write(AdditionalData);

            // Write player properties
            writer.Write(Players.Count);
            foreach (var playerProperty in Players)
            {
                writer.Write(playerProperty.Key);

                // Write the member info
                playerProperty.Value.ToBinaryWriter(writer);
            }

            // Write teams info
            writer.Write(Teams.Count);
            foreach (var team in Teams)
            {
                writer.Write(team.Key);

                // Write team data
                team.Value.ToBinaryWriter(writer);
            }

            // Write controls
            writer.Write(Controls.Count);
            foreach (var control in Controls)
            {
                writer.Write(control.Key);

                control.Value.ToBinaryWriter(writer);
            }

            // Other settings
            writer.Write(EnableManualStart);
            writer.Write(EnableReadySystem);
            writer.Write(EnableTeamSwitching);

        }

        public override void FromBinaryReader(EndianBinaryReader reader)
        {
            LobbyState = (LobbyState) reader.ReadInt32();
            StatusText = reader.ReadString();
            GameMaster = reader.ReadString();
            CurrentUserUsername = reader.ReadString();

            LobbyId = reader.ReadInt32();
            LobbyName = reader.ReadString();
            LobbyProperties = reader.ReadDictionary();
            MaxPlayers = reader.ReadInt32();

            // Read additional data
            var size = reader.ReadInt32();
            if (size > 0) AdditionalData = reader.ReadBytes(size);

            // Clear, in case we're reusing the object
            Players.Clear();

            // Read player properties
            var playerCount = reader.ReadInt32();

            for (var i = 0; i < playerCount; i++)
            {
                var data = CreateLobbyMemberData();
                var username = reader.ReadString();
                data.FromBinaryReader(reader);

                Players.Add(username, data);
            }

            // Read teams
            Teams.Clear();
            var teamsCount = reader.ReadInt32();
            for (int i = 0; i < teamsCount; i++)
            {
                var teamKey = reader.ReadString();
                var teamData = CreateTeamData();
                teamData.FromBinaryReader(reader);
                Teams.Add(teamKey, teamData);
            }

            // Read controls
            Controls.Clear();
            var controlsCount = reader.ReadInt32();
            for (int i = 0; i < controlsCount; i++)
            {
                var controlKey = reader.ReadString();
                var controlData = CreateLobbyPropertyData();
                controlData.FromBinaryReader(reader);
                Controls.Add(controlKey, controlData);
            }

            // Other settings
            EnableManualStart = reader.ReadBoolean();
            EnableReadySystem = reader.ReadBoolean();
            EnableTeamSwitching = reader.ReadBoolean();
        }

        protected virtual LobbyMemberData CreateLobbyMemberData()
        {
            return new LobbyMemberData();
        }

        protected virtual LobbyTeamData CreateTeamData()
        {
            return new LobbyTeamData();
        }

        protected virtual LobbyPropertyData CreateLobbyPropertyData()
        {
            return new LobbyPropertyData();
        }
    }
}