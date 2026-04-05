package profile

import "fmt"

func showProfile(profile string) (err error) {
	config, err := LoadProfile(profile)
	if err != nil {
		return err
	}

	fmt.Printf("Profile [%s]:\n", profile)
	fmt.Printf("  Server URL: %s\n", config.URL)

	return nil
}
