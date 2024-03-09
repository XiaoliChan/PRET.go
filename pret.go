package main

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func pjl(target string) {
	info_command := "@PJL INFO ID"
	cmd := create_PJLCMD("@PJL USTATUSOFF", false)
	cmd_info := create_PJLCMD(info_command, true)

	conn := connectTCP(target)
	if conn != nil {
		_, err := conn.Write([]byte(cmd))
		if err == nil {
			_, _ = GetNetResponse(conn)
			_, err := conn.Write([]byte(cmd_info))
			if err == nil {
				result, _ := GetNetResponse(conn)
				if len(result) == 0 {
					fmt.Println("Non-PJL Device")
				} else {
					fmt.Println(process_result_PJL(result, info_command))
				}
			}
		}
	}
}

func ps(target string) {
	iohack := true
	// check iohack
	stage_1 := "(x1) = (x2) == << /DoPrintErrors false >> setsystemparams"
	stage_2 := "product print"
	// iohack is true here
	cmd, _ := create_PSCMD(stage_1, iohack)
	conn := connectTCP(target)
	if conn != nil {
		_, err := conn.Write([]byte(cmd))
		if err == nil {
			//res, _ := ReadBytes(conn)
			res, _ := GetNetResponse(conn)
			if len(res) != 0 {
				if strings.Contains(string(res), "x1") {
					iohack = false
				}
				if !strings.Contains(string(res), "x1") || !strings.Contains(string(res), "x2") {
					fmt.Println("iohack is true")
				}
				cmd, token := create_PSCMD(stage_2, iohack)
				_, err := conn.Write([]byte(cmd))
				if err == nil {
					res, _ = GetNetResponse(conn)
					fmt.Println(process_result_PS(res, token))
				}
			}
		}
	}
}

func pcl(target string) {
	cmd := create_PCLCMD("*s1M")
	conn := connectTCP(target)
	if conn != nil {
		_, err := conn.Write([]byte(cmd))
		if err == nil {
			res, _ := GetNetResponse(conn)
			if len(res) != 0 {
				if strings.Contains(string(res), "PCL") {
					fmt.Println(process_result_PCL(res))
				} else {
					fmt.Println("Non-PCL Device")
				}
			}
		}
	}
}

func create_PSCMD(cmd string, iohack bool) (string, string) {
	send := fmt.Sprintf("{%s} stopped", cmd)
	token := RawCommand["DELIMITER"] + RawCommand["RAND_NUM"]
	footer := "\n(" + token + "\\n) print flush\n"
	if !iohack {
		return RawCommand["UEL"] + RawCommand["PS_HEADER"] + send + footer, token
	} else {
		return RawCommand["UEL"] + RawCommand["PS_HEADER"] + RawCommand["PS_IOHACK"] + send + footer, token
	}
}

func create_PJLCMD(cmd string, wait bool) string {
	enable_status := false
	var (
		status string
		footer string
	)
	token := RawCommand["DELIMITER"] + RawCommand["RAND_NUM"]
	if wait {
		if enable_status {
			status = "@PJL INFO STATUS" + RawCommand["EOL"]
		}
		footer = "@PJL ECHO " + token + RawCommand["EOL"] + RawCommand["EOL"]
	} else {
		status = ""
		footer = ""
	}

	return RawCommand["UEL"] + cmd + RawCommand["EOL"] + status + footer + RawCommand["UEL"]
}

func create_PCLCMD(cmd string) string {
	token := RawCommand["RAND_NUM_PCL"]
	footer := RawCommand["ESC"] + "*s" + token + "X"
	return RawCommand["UEL"] + RawCommand["PCL_HEADER"] + cmd + footer + RawCommand["UEL"]
}

func process_result_PJL(result []byte, info_command string) string {
	info_command_Hex := hex.EncodeToString([]byte(info_command)) + "0d0a"
	footer_Hex := "0d0a0c40504a4c204543484f2044454c494d49544552" // 0d0a0c + hex("@PJL ECHO DELIMITER")
	result_Hex := hex.EncodeToString(result)
	if strings.Contains(result_Hex, info_command_Hex) {
		var result string
		result_ := strings.Split(result_Hex, info_command_Hex)
		if len(result_) > 1 {
			result = result_[1]
		} else {
			result = result_[0]
		}
		result_Processed := strings.Split(result, footer_Hex)[0]
		bytes, _ := hex.DecodeString(result_Processed)
		return string(bytes)
	} else {
		return "Failed"
	}
}

func process_result_PCL(result []byte) string {
	result_Hex := hex.EncodeToString(result)
	if strings.Contains(result_Hex, "50434c0d0a") {
		result_Hex = strings.ReplaceAll(result_Hex, "50434c0d0a", "")
		result_Hex = strings.ReplaceAll(result_Hex, "0d0a0c", "")
		bytes, _ := hex.DecodeString(result_Hex)
		return string(bytes)
	} else {
		return "Failed"
	}
}

func process_result_PS(result []byte, token string) string {
	result_Hex := hex.EncodeToString(result)
	if strings.Contains(result_Hex, "78310d0a78320a") {
		token_hex := hex.EncodeToString([]byte(token)) + "0a"
		if strings.Contains(result_Hex, token_hex) {
			var result string
			result_Hex_ := strings.Split(result_Hex, token_hex)
			if len(result_Hex_) > 1 {
				result = result_Hex_[1]
			} else {
				result = result_Hex_[0]
			}
			result_Processed, _ := hex.DecodeString(strings.ReplaceAll(result, "0d0a", ""))
			return string(result_Processed)
		} else {
			return "Failed"
		}
	} else {
		token_hex := hex.EncodeToString([]byte(token))
		if strings.Contains(result_Hex, token_hex) {
			var result string
			result_Hex_ := strings.Split(result_Hex, token_hex+"0a")
			if len(result_Hex_) > 1 {
				result = result_Hex_[1]
			} else {
				result = result_Hex_[0]
			}
			result_Processed, _ := hex.DecodeString(strings.ReplaceAll(result, token_hex+"0d0a", ""))
			return string(result_Processed)
		} else {
			return "Failed"
		}
	}
}
