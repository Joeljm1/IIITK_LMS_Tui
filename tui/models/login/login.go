package login

import (
	"errors"
	"strings"

	"github.com/Joeljm1/IIITKlmsTui/internal/client"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zalando/go-keyring"
)

const serviceName = "lmsTui"

var (
	errLogin       = LoginErr{errors.New("failed to fetch login details")}
	errValidDetail = LoginValidErr{errors.New("details were not valid")}
)

// Get username from os kernal keyring
func (m Model) getUsernameFromKeyRing() (string, error) {
	uname, err := keyring.Get(serviceName, "username")
	if err != nil {
		return "", err
	}
	return uname, nil
}

// Send uname to update if present else error
// Used in init fn
func (m Model) sendUsername() tea.Msg {
	uname, err := m.getUsernameFromKeyRing()
	if err != nil {
		return errLogin
	}
	return UName(uname)
}

// Get psswd from os kernal keyring
func (m Model) getPsswdFromKeyRing() (string, error) {
	psswd, err := keyring.Get(serviceName, "password")
	if err != nil {
		return "", err
	}
	return psswd, nil
}

// Send psswd to update if present else error
// Used in init fn
func (m Model) sendPasswd() tea.Msg {
	psswd, err := m.getPsswdFromKeyRing()
	if err != nil {
		return errLogin
	}
	return Psswd(psswd)
}

func (m Model) setDetailsToKeyring(username, password string) error {
	err := keyring.Set(serviceName, "username", username)
	if err != nil {
		return err
	}
	err = keyring.Set(serviceName, "password", password)
	if err != nil {
		return err
	}
	return nil
}

func (m Model) LoginComplete() tea.Msg {
	return LoginComplete("Complete")
}

func (m Model) deleteLoginDet() error {
	err := keyring.DeleteAll(serviceName)
	return err
}

// Return err if uname and psswd not valid else Valid(true).
// If valid sets to uname nad psswd to keyring to
func (m Model) validateDetails(username, password string) tea.Cmd {
	return func() tea.Msg {
		if str := strings.Trim(username, " "); len(str) < 11 {
			return errValidDetail
		}
		if str := strings.Trim(password, " "); len(str) < 8 {
			return errValidDetail
		}
		_, err := client.NewClient(username, password)
		if err != nil {
			return errValidDetail
		}
		err = m.setDetailsToKeyring(username, password)
		if err != nil {
			return errValidDetail
		}
		return Valid(true)
	}
}
