package api

var authType = []string{
	"basic-auth", // User: See: https://github.com/coderbydesign/identity-schemas/blob/add-validator/3scale/schema.json#L51
	"jwt-auth",   // User: See: https://github.com/coderbydesign/identity-schemas/blob/add-validator/3scale/schema.json#L51
	"cert-auth",  // System: See: https://github.com/coderbydesign/identity-schemas/blob/add-validator/3scale/schema.json#L113
	"uhc-auth",   // System: See: https://github.com/coderbydesign/identity-schemas/blob/add-validator/3scale/schema.json#L113
}
