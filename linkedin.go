/*

TLDR;
goto https://www.linkedin.com/oauth/v2/authorization?response_type=code&state=123456789&scope=w_member_social%2Cr_liteprofile&client_id=77kgucrnlc71n7&redirect_uri=https%3A%2F%2Fgoogle.com
check the URL you will find code=<code>
get the code and update the code variable


Post Account setup we have below information from the provider End
Static information
Client ID: 77kgucrnlc71n7
Client Secret: rwhpdHzwM5LNzK3Z
Redirect URL: https://google.com
URL Encoded: https%3A%2F%2Fgoogle.com
State: w_member_social%2Cr_liteprofile

User will need to do setup to get the code details. This will redirect user to linkedin form and
user must allow provider to access the information and then we can proceed with below code.

Authorization Endpoint:-
https://www.linkedin.com/oauth/v2/authorization?response_type=code&state=123456789&scope=w_member_social%2Cr_liteprofile&client_id=77kgucrnlc71n7&redirect_uri=https%3A%2F%2Fgoogle.com

Response:
https://www.google.com/?code=AQSZZelPf-m1v-rrrBksNiFBlH4eQCwZdas5mv3f1xHTa5VoGlcCUV6HdK5ui6WKh5v1GnaMwMuIY-9W6q6T9MYNQEWebfny0k1iMmicD-RNmrK-ylJO_LtDYcWpdj9nIuIUUSRcPNHseze7Bgpw4eb6pttaLR5_XVW-S-RSxPQtuiV25qcSyAvRXzLtvBVRLon5JysjmeArVPnRQqc&state=123456789

Fliter:
code=AQSZZelPf-m1v-rrrBksNiFBlH4eQCwZdas5mv3f1xHTa5VoGlcCUV6HdK5ui6WKh5v1GnaMwMuIY-9W6q6T9MYNQEWebfny0k1iMmicD-RNmrK-ylJO_LtDYcWpdj9nIuIUUSRcPNHseze7Bgpw4eb6pttaLR5_XVW-S-RSxPQtuiV25qcSyAvRXzLtvBVRLon5JysjmeArVPnRQqc

curl -ik -X POST https://www.linkedin.com/oauth/v2/accessToken
			-d grant_type=authorization_code
			-d code=AQSZZelPf-m1v-rrrBksNiFBlH4eQCwZdas5mv3f1xHTa5VoGlcCUV6HdK5ui6WKh5v1GnaMwMuIY-9W6q6T9MYNQEWebfny0k1iMmicD-RNmrK-ylJO_LtDYcWpdj9nIuIUUSRcPNHseze7Bgpw4eb6pttaLR5_XVW-S-RSxPQtuiV25qcSyAvRXzLtvBVRLon5JysjmeArVPnRQqc
			-d redirect_uri=https%3A%2F%2Fgoogle.com
			-d client_id=77kgucrnlc71n7
			-d client_secret=rwhpdHzwM5LNzK3Z



*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	code          = ""
	client_id     = "77kgucrnlc71n7"
	client_secret = "rwhpdHzwM5LNzK3Z"
	redirect_url  = "https://google.com"
	timeout       = time.Duration(5 * time.Second)
)

var client = http.Client{
	Timeout: timeout,
}

func getToken() (string, error) {

	// if token exist in the file then exit
	saved_token, _ := ioutil.ReadFile("linkedinToken.text")
	if len(saved_token) > 0 {
		token := map[string]string{}
		_ = json.Unmarshal(saved_token, &token)
		return fmt.Sprint(token["access_token"]), nil
	}
	// client decleration

	body := url.Values{}
	body.Set("grant_type", "authorization_code")
	body.Set("code", code)
	body.Set("redirect_uri", redirect_url)
	body.Set("client_id", client_id)
	body.Set("client_secret", client_secret)

	uri := "https://www.linkedin.com/oauth/v2/accessToken"
	u, _ := url.ParseRequestURI(uri)
	urlStr := u.String()

	req, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(body.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		log.Fatalln("System error while creating request\nErr: ", err)
	}

	// calling action
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln("System error while getting access token\nErr: ", err)
	}

	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln("System Error while reading response that contains access token/nErr: ", err)
	}

	token := map[string]string{}
	_ = json.Unmarshal(resp_body, &token)

	string_body := string(resp_body)

	if len(fmt.Sprintf(token["access_token"])) == 0 {
		return "", fmt.Errorf(string_body)
	}

	// saving token in a file
	byte_body := []byte(string_body)
	// the WriteFile method returns an error if unsuccessful
	_ = ioutil.WriteFile("linkedinToken.text", byte_body, 0777)

	return fmt.Sprint(token["access_token"]), nil
}

func getUserDetails(token string) (map[string]string, error) {
	url := "https://api.linkedin.com/v2/me"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	if err != nil {
		log.Fatalln("System error while creating request\nErr: ", err)
	}

	// calling action
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln("System error while getting user details from linkedin\nErr: ", err)
	}

	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln("System Error while reading response/nErr: ", err)
	}

	details := map[string]string{}

	_ = json.Unmarshal(resp_body, &details)

	string_body := string(resp_body)

	if len(fmt.Sprintf(details["localizedLastName"])) == 0 {
		return nil, fmt.Errorf(string_body)
	}
	return details, nil
}

func linkedInPost(message string, token string, id string) error {

	url := "https://api.linkedin.com/v2/posts"

	payload := fmt.Sprintf(`{
		"author": "urn:li:person:%s",
		"commentary": "%s",
		"visibility": "PUBLIC",
		"distribution": {
		  "feedDistribution": "MAIN_FEED",
		  "targetEntities": [],
		  "thirdPartyDistributionChannels": []
		},
		"lifecycleState": "PUBLISHED",
		"isReshareDisabledByAuthor": false
	  }`, id, message)

	fmt.Println(payload)
	bytes_payload := []byte(payload)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bytes_payload))

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Restli-Protocol-Version", "2.0.0")

	if err != nil {
		log.Fatalln("System error while creating request\nErr: ", err)
	}

	// calling action
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln("System error while getting user details from linkedin\nErr: ", err)
	}

	defer resp.Body.Close()

	if err != nil {
		log.Fatalln("System Error while reading response\nErr: ", err)
	}

	if resp.Status == "201" {
		return fmt.Errorf("Unable to create the post")
	}

	return nil
}

func main() {
	token, err := getToken()

	if err != nil {
		fmt.Print(err)
		return
	}

	resp, err := getUserDetails(token)
	if err != nil {
		fmt.Print("Err: ", err)

	}
	message := "This message is from a automated posting program, please ignore till the owner comes and delete this post manually."
	err = linkedInPost(message, token, resp["id"])
}
