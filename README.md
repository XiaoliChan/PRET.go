# PRET.go
Printer Exploitation Toolkit with golang

## Example

### Show printer info
```
package main

// Test-with: HP LaserJet M402dw

func main() {
	target := "192.168.1.10:9100"
	ps(target)
	//pjl(target)
	//pcl(target)
}
```

Output  
![image](https://github.com/XiaoliChan/PRET.go/assets/30458572/82eafed2-1fde-40a4-b7af-1a9c86d2d08d)
