package conf

import "time"

const ProjectName = "hc-management"

var RunMode = "normal" // normal, test

// RedisPassword ...
const RedisPassword string = ""

// Version of Server
const Version string = "/v1"

// OkCode ...
const OkCode int = 1

// InterErrorCode ...
const InterErrorCode int = 2

// BadParamsCode ...
const BadParamsCode int = 3

const MinEntropyBits = 60

// jwt access_token lifetime
const AccessToken_Lifetime = time.Second * 86400 * 7

// jwt refresh_token lifetime
const RefreshToken_Lifetime = time.Second * 86400 * 14

// jwt refresh_token lifetime
const BanIpTime = time.Minute * 15

const SendCodeTime = time.Minute * 5

const SetOTPCodeValidatation = time.Minute * 30

const SendCodeTimeDuration = time.Minute * 3
