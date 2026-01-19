package zotero

func (z *ZoteroClient) FetchCollections() ([]Collection, error) {
	url := z.buildURL("collections")
	return fetch[Collection](z, url)
}
