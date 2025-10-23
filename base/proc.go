package base

import (
	"strings"
)

var Proc = make(map[string][]byte)

func InitProc(ip string) {
	Proc["TCP"] = []byte("GET / HTTP/1.1\r\nHost: " + ip + "\r\nConnection: close\r\n\r\n")
	Proc["RDP"] = []byte{
		0x03, 0x00, 0x00, 0x13, // TPKT Header: Version 3, Reserved, Length (19 bytes)
		0x0e, 0xe0, 0x00, 0x00, // X.224 Header: Length (14 bytes), Type (CR)
		0x00, 0x00, 0x00, // Dest Ref, Src Ref, Class

		// --- RDP 协商请求 ---
		0x01, 0x00, // Type: RDP_NEG_REQ (0x01)
		0x08, 0x00, // Length (8 bytes)
		0x03, 0x00, 0x00, 0x00, // Requested Protocols: PROTOCOL_HYBRID (3)
	}
	Proc["REDIS"] = []byte("PING\r\n")
}

func GetInfo(reqProc string, resp []byte) string {
	if len(resp) == 0 {
		return ""
	}
	if reqProc == "TCP" {
		if len(resp) < 3 {
			return ""
		}
		sresp := string(resp)
		// fmt.Println(sresp)
		if sresp[:3] == "SSH" {
			return "SSH"
		}
		if sresp[:4] == "HTTP" {
			return "HTTP "
		}
		if sresp[:3] == "220" {
			return sresp[4:20]
		}
		if sresp[:3] == "+OK" {

			return sresp[4:20]
		}
		if sresp[:3] == "RFB" {
			return "VNC"
		}
		if len(sresp) >= 17 && sresp[:17] == "* OK [CAPABILITY " {
			return sresp[17:22]
		} else if sresp[:4] == "* OK" {
			return sresp[5:15]
		}

		if strings.Contains(sresp, "mysql_native_password") || strings.Contains(sresp, "caching_sha2_password") {
			return "MYSQL"
		}

		return "Lain Don't Know"
	}
	if reqProc == "REDIS" {
		if string(resp)[:5] == "+PONG" {
			return "REDIS"
		}
	}
	// fmt.Println(resp)
	if reqProc == "RDP" && len(resp) > 8 {
		if resp[0] == 0x03 && resp[1] == 0x00 && resp[2] == 0 && resp[5] == 0xd0 || resp[5] == 0xf0 {
			return "RDP"
		}
	}
	return ""
}
