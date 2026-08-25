package cli

func LoginCopy() string {
	return "Connect your wallet in the desktop or web app, then sign the bind message. PIT never asks for your private key."
}

func RefuseLoginSecret(line string) error {
	return RefusePrint(line)
}
