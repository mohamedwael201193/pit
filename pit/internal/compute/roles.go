package compute

func Committee(publicMarket, privateBook []byte) (map[Role][]byte, error) {
	out := map[Role][]byte{}
	for _, r := range []Role{Researcher, Challenger, Risk} {
		b, err := Envelope(r, publicMarket, privateBook)
		if err != nil {
			return nil, err
		}
		out[r] = b
	}
	return out, nil
}

func IndependenceNote() Independence {
	return EnvelopeOnly
}
