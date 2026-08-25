module github.com/mohamedwael201193/pit/sealer

go 1.24.0

require (
	github.com/0gfoundation/0g-pc-e2ee/client v0.0.0
	github.com/0gfoundation/0g-pc-e2ee/protocol v0.0.0
)

require (
	github.com/cloudflare/circl v1.6.4 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.0 // indirect
	github.com/gowebpki/jcs v1.0.1 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
)

replace github.com/0gfoundation/0g-pc-e2ee/client => ./third_party/client

replace github.com/0gfoundation/0g-pc-e2ee/protocol => ./third_party/protocol
