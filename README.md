# go_oura

A way to pull the data for your Oura ring using Go.

Uses the v2 of the [Oura Ring API](https://cloud.ouraring.com/v2/docs).

## Features

- Oura Ring v2 API
  - [Daily Activity](https://cloud.ouraring.com/v2/docs#tag/Daily-Activity-Routes)
  - [Daily Readiness](https://cloud.ouraring.com/v2/docs#tag/Daily-Readiness-Routes)
  - [Daily Sleep](https://cloud.ouraring.com/v2/docs#tag/Daily-Sleep-Routes)
  - [Daily Spo2](https://cloud.ouraring.com/v2/docs#tag/Daily-Spo2-Routes)
  - [Daily Stress](https://cloud.ouraring.com/v2/docs#tag/Daily-Stress-Routes)
  - [Tags](https://cloud.ouraring.com/v2/docs#tag/Enhanced-Tag-Routes)
  - [Heart Rate](https://cloud.ouraring.com/v2/docs#tag/Heart-Rate-Routes)
  - [Personal Info](https://cloud.ouraring.com/v2/docs#tag/Personal-Info-Routes)
  - [Rest Mode](https://cloud.ouraring.com/v2/docs#tag/Rest-Mode-Period-Routes)
  - [Ring Configuration](https://cloud.ouraring.com/v2/docs#tag/Ring-Configuration-Routes)
  - [Session](https://cloud.ouraring.com/v2/docs#tag/Session-Routes)
  - [Sleep](https://cloud.ouraring.com/v2/docs#operation/Multiple_sleep_Documents_v2_usercollection_sleep_get)
  - [Sleep Time](https://cloud.ouraring.com/v2/docs#tag/Sleep-Time-Routes)
  - [Workout](https://cloud.ouraring.com/v2/docs#tag/Workout-Routes)

## What's Missing

- OAuth2 Authentication &rarr; So far go_oura requires a personal token to authenticate against the Oura Ring API. OAuth2 protocol may be added in the future. However, as of now this is geared towards personal usage, and not for creating an Oura registered application.
- [Webhooks](https://cloud.ouraring.com/v2/docs#tag/Webhook-Subscription-Routes) &rarr; Would be used mainly for an Oura registered application. At this time has not been implemented.

## Install

```bash
go get github.com/austinmoody/go_oura
```

## Examples

The [examples](examples) directory contains an example for each implemented type.  Each example makes a couple calls to the Oura Ring API and displays data.

To run the examples you'll need to have an environment variable called ```OURA_ACCESS_TOKEN``` setup to contain your personal access token.

## Usage

The first thing to do is to get a Client.  

```go
client := go_oura.NewClient("your-access-token")
```

Next we can pull data for different types of data offered up by Oura Ring.

Activities for the last 24 hours:

```go
activities, err := client.GetActivities(time.Now().Add(-24 * time.Hour), time.Now(), nil)
if err != nil {
    fmt.Printf("Error getting activities: %v", err)
    return
}
```

Unless there is an error you'll have an [Activities](activity.go) type.