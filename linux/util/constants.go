package util

import "time"

const (
	Status            string        = "status"
	Hostname          string        = "hostname"
	Err               string        = "err"
	Data              string        = "data"
	CredentialProfile string        = "credential_profile"
	Username          string        = "username"
	Password          string        = "password"
	SSHTimout         time.Duration = 30
	TCP               string        = "tcp"
	DiscoveryProfile  string        = "discovery_profile"
	IP                string        = "ip"
	Port              string        = "port"
	Empty             string        = ""
	NewLine           string        = "\n"
	Space             string        = " "
	NotDiscovered     string        = "Device not discovered!"
	Discovered        string        = "Device discovered successfully!"
	Success           string        = "Success"
	Fail              string        = "Fail"
	Colon             string        = ":"
	SystemProcess     string        = "system.process"
	ProcessPID        string        = "system.process.pid"
	ProcessCPU        string        = "system.process.cpu"
	ProcessMemory     string        = "system.process.memory"
	ProcessUser       string        = "system.process.user"
	ProcessCommand    string        = "system.process.command"
)
