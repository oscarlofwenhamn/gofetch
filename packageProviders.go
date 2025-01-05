package main

import (
	"fmt"
	"os/exec"
)

func GetSnapPackageCount() (string, error) {
	snapList := exec.Command("snap", "list")
	out, err := snapList.Output()
	if err != nil {
		return "", err
	}

	var count int
	for _, b := range out {
		if b == '\n' {
			count++
		}
	}
	count--
	return fmt.Sprintf("%d (snap)", count), nil
}

func GetDpkgPackageCount() (string, error) {
	dpkgList := exec.Command("dpkg-query", "-f", "'.\n'", "-W")
	out, err := dpkgList.Output()
	if err != nil {
		return "", err
	}

	var count int
	for _, b := range out {
		if b == '\n' {
			count++
		}
	}

	return fmt.Sprintf("%d (dpkg)", count), nil
}
