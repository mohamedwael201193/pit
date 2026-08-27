package cli

import (
	"github.com/mohamedwael201193/pit/internal/config"
	pitexec "github.com/mohamedwael201193/pit/internal/exec"
	"github.com/mohamedwael201193/pit/internal/hl"
)

func LiveAsset(network, coin string) (hl.BookSnapshot, error) {
	net, err := config.ParseNetwork(network)
	if err != nil {
		return hl.BookSnapshot{}, err
	}
	book, err := hl.New(config.For(net)).PublicBook(coin)
	if err != nil {
		return hl.BookSnapshot{}, err
	}
	if err := pitexec.NeedAssetIndex(true, book.Asset); err != nil {
		return hl.BookSnapshot{}, err
	}
	return book, nil
}
