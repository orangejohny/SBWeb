// Copyright Dmitry Kargashin <dkargashin3@gmail.com>
// License can be found in LICENSE file.

/*
  Package api implements handling of different URL addresses of application.
	It uses data types provided by model package.

	Data types

	Errors

	Package api uses not only default HTTP status codes but own error type
	that describes happened error.

	API Error:
		description      short description of error
		message          what client should do for error resolving
		error code       unique error code

	User

	Names of fields of JSON object which will be returned:
		id               <int64>
		first_name       <string>  [UTF letter]
		last_name        <string>  [UTF letter]
		email            <string>
		tel_number       <string>  [digits 1-9]
		about            <string>  [ASCII]
		reg_time         <string>


	HTTP parameters

	Interface

	Read and search multiple ads

	"root" is domain part of server (i.e. http://example.com)

	"root/ads" address:
		method                 GET
		allowed parameters:
			query                search query; return only ads which contatins query in title of ad
			limit                maximum number of ads which will be returned
			offset               number of the first ad that will be returned
		return result:
			status 200 and JSON array of ads
*/

package api
