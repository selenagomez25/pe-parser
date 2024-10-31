package main

import (
	"debug/pe"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: pe-parser <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	file, err := pe.Open(filename)
	if err != nil {
		fmt.Printf("[!] error: %v\n", strings.ToLower(err.Error()))
		os.Exit(1)
	}
	defer file.Close()

	fmt.Printf("[+] file: %s\n\n", filename)

	machineType := func(machine uint16) string {
		switch machine {
		case 0x014c:
			return "x86"
		case 0x8664:
			return "x64"
		case 0x01c4:
			return "ARM"
		case 0xaa64:
			return "ARM64"
		default:
			return fmt.Sprintf("unknown (0x%x)", machine)
		}
	}

	characteristics := func(char uint16) string {
		var chars []string
		if char&0x0020 != 0 {
			chars = append(chars, "large address aware")
		}
		if char&0x0100 != 0 {
			chars = append(chars, "32 bit")
		}
		if char&0x2000 != 0 {
			chars = append(chars, "dll")
		}
		if len(chars) == 0 {
			return "none"
		}
		return strings.Join(chars, ", ")
	}

	fmt.Printf("[+] machine: %s\n", machineType(file.FileHeader.Machine))
	fmt.Printf("[+] sections: %d\n", file.FileHeader.NumberOfSections)
	fmt.Printf("[+] timestamp: %v\n", file.FileHeader.TimeDateStamp)
	fmt.Printf("[+] characteristics: %s\n", characteristics(file.FileHeader.Characteristics))

	switch oh := file.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		fmt.Printf("[+] architecture: 32-bit\n")
		fmt.Printf("[+] imagebase: 0x%x\n", oh.ImageBase)
		fmt.Printf("[+] entrypoint: 0x%x\n", oh.AddressOfEntryPoint)
	case *pe.OptionalHeader64:
		fmt.Printf("[+] architecture: 64-bit\n")
		fmt.Printf("[+] imagebase: 0x%x\n", oh.ImageBase)
		fmt.Printf("[+] entrypoint: 0x%x\n", oh.AddressOfEntryPoint)
	}

	fmt.Println("\n[+] sections:")
	for _, section := range file.Sections {
		fmt.Printf("  %-9s  size: %-8d  addr: 0x%-8x\n",
			section.Name, section.Size, section.VirtualAddress)
	}
}
