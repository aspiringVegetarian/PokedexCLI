package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (cl *Client) ListLocationAreas(pageURL *string) (LocationAreaResp, error) {
	fullURL := baseURL + "/location-area?offset=0&limit=20"
	if pageURL != nil {
		fullURL = *pageURL
	}

	// check the cache

	if data, exists := cl.cache.Get(fullURL); exists {
		fmt.Println("Cache hit!")
		locationAreaResp := LocationAreaResp{}

		err := json.Unmarshal(data, &locationAreaResp)
		if err != nil {
			return LocationAreaResp{}, err
		}

		return locationAreaResp, nil
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return LocationAreaResp{}, err
	}

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return LocationAreaResp{}, err
	}

	if resp.StatusCode > 399 {
		return LocationAreaResp{}, fmt.Errorf("bad status code: %v", resp.StatusCode)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreaResp{}, err
	}

	// add to cache
	cl.cache.Add(fullURL, data)

	locationAreaResp := LocationAreaResp{}

	err = json.Unmarshal(data, &locationAreaResp)
	if err != nil {
		return LocationAreaResp{}, err
	}

	return locationAreaResp, nil
}

func (cl *Client) ExploreLocationArea(specificLocation *string) (SpecificLocationAreaResp, error) {

	fullURL := baseURL + "/location-area/" + *specificLocation

	// check the cache

	if data, exists := cl.cache.Get(fullURL); exists {
		fmt.Println("Cache hit!")
		specificLocationAreaResp := SpecificLocationAreaResp{}

		err := json.Unmarshal(data, &specificLocationAreaResp)
		if err != nil {
			return SpecificLocationAreaResp{}, err
		}

		return specificLocationAreaResp, nil
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return SpecificLocationAreaResp{}, err
	}

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return SpecificLocationAreaResp{}, err
	}

	if resp.StatusCode > 399 {
		return SpecificLocationAreaResp{}, fmt.Errorf("bad status code: %v", resp.StatusCode)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return SpecificLocationAreaResp{}, err
	}

	// add to cache
	cl.cache.Add(fullURL, data)

	specificLocationAreaResp := SpecificLocationAreaResp{}

	err = json.Unmarshal(data, &specificLocationAreaResp)
	if err != nil {
		return SpecificLocationAreaResp{}, err
	}

	return specificLocationAreaResp, nil
}
