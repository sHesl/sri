# sri

SRI is a CLI tool for generating Sub-Resource Integrity digests for Javascript files. 

`go get github.com/sHesl/sri`

## Usage

`sri .` to generate SRI digests for all files in current directory .   
`sri jquery-3.3.1.min.js` to generate SRI digests for a single, local file .   
`sri https://code.jquery.com/jquery-3.3.1.min.js` to generate SRI digests for a single, hosted file .   


## Example Output
SRI produces a JSON file with digests for sha256/384/512, as well as the relevant script tag with integrity attribute.  
```
{
	"jquery-3.3.1.min.js": {
		"sha256": {
			"digest": "sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8=",
			"tag": "<script src='https://code.jquery.com/jquery-3.3.1.min.js' integrity='sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8='></script>"
		},
		"sha384": {
			"digest": "sha384-tsQFqpEReu7ZLhBV2VZlAu7zcOV+rXbYlF2cqB8txI/8aZajjp4Bqd+V6D5IgvKT",
			"tag": "<script src='https://code.jquery.com/jquery-3.3.1.min.js' integrity='sha384-tsQFqpEReu7ZLhBV2VZlAu7zcOV+rXbYlF2cqB8txI/8aZajjp4Bqd+V6D5IgvKT'></script>"
		},
		"sha512": {
			"digest": "sha512-+NqPlbbtM1QqiK8ZAo4Yrj2c4lNQoGv8P79DPtKzj++l5jnN39rHA/xsqn8zE9l0uSoxaCdrOgFs6yjyfbBxSg==",
			"tag": "<script src='https://code.jquery.com/jquery-3.3.1.min.js' integrity='sha512-+NqPlbbtM1QqiK8ZAo4Yrj2c4lNQoGv8P79DPtKzj++l5jnN39rHA/xsqn8zE9l0uSoxaCdrOgFs6yjyfbBxSg=='></script>"
		}
	}
}
```
