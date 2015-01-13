[![Build Status](https://travis-ci.org/sqdk/samurai.svg)](https://travis-ci.org/sqdk/samurai)

samurai - a string tokenization library for Go
=======

# About
Samurai takes data in a regular format (like logfiles) and tokenizes it according to a predefined pattern. The tokenization is done by gradually splitting (or slicing, hence the name) the data until you get the subcomponents that you want. 

I made samurai as a functionally-equivalent alternative to grok.

In its current state, i have been able to tokenize approx. 1.6 million lines (160 mb) of apache logs in 15-20 seconds (about 0.01 ms pr line) single-threaded on a Core i7-4500U@2.40GHz. Current experiments with goroutines created a massive memoryleak, but the code should in theory be threadable. 

If you want to test it out yourself, big log files can be found here: http://ita.ee.lbl.gov/html/contrib/NASA-HTTP.html

## Pattern syntax
Lets start with a simple example. Our input data is a collection of semi-colon seperated lines of data:

```
	John;21;555-32132-11
	Matt;32;555-11231-11
	Chris;32;555-32211-32
```

We can tokenize this data by using the following pattern:
```
	pattern := ";(firstName,age,tel)"
	//; is the char or string to split the data with
	//The names inside the parenthesis indicate the given names of the values at the
	//relative index in the resulting array after splitting.
```

Which is the equivalent to these string operations:
```
	inputString := "John;21;555-32132-11"
	subComponents := strings.Split(inputString, ";")
	values := make(map[string]string)

	values["firstName"] = values[0]
	values["age"] = values[1]
	values["tel"] = values[2]

	return values
```


Patterns can also be nested. Lets extend the previous data to include last names:
```
	John Johnson;21;555-32132-11
	Matt Mattson;32;555-11231-11
	James Jameson;32;555-32211-32
```

If we used the last pattern, we would end up with both the first and last name as one value. If we the first and last name as seperate values, we can insert a nested pattern that does this:
```
	pattern := ";( (firstName,lastName),age,tel)"
```

# Examples
## Tokenizing apache logs
A standard apache log line looks like this:
```
	127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326
```

This pattern can be used to tokenize the data:
```
	apacheLogPattern := "[( (ip,nil,user),](date,\"(nil, (method,url,httpver), (nil,httpcode,reqsize))))"
```
Note that the pattern shows the format of the data in a much simpler way that the equivalent regex pattern. You can even see some kind of context as to what data it is meant to extract.

Data can typically be split in multiple ways. This is an alternative pattern that is slightly faster because of reduced nesting of patterns:
```
	exampleWithoutDelimiter := " (ip,nil,user,[(nil,date),](tz),\"(nil,method),url,\"(httpver), code, size)"
```
The only difference is you get date and timezone as separate values because the first split is a space.
