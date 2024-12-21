package provider

type ProviderVonage struct {
	apiKey    string
	apiSecret string
}

func (p *ProviderVonage) Send(opts SMSProviderSendOpts) error {
	// TODO: Implement this
	return nil
}
