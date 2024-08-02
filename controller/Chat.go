func makeRequest(client *resty.Client, url, token, prompt string) (*resty.Response, error) {
	response, err := client.R().
		SetHeader("Authorization", token).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"inputs": prompt}).
		Post(url)

	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		return response, errors.New("unexpected status code: " + strconv.Itoa(response.StatusCode()))
	}

	return response, nil
}
