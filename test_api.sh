#!/bin/bash
api_url="http://localhost:1323"

test_full(){
    test_signup
}
test_signup(){
    echo "post SignUpData"
    curl -v -d "@test/SignUpData.json" -H "Content-Type: application/json" POST ${api_url}/auth/signup
}
test_confirm(){
    echo "get /signup/confirm/:token"
    curl -v GET ${api_url}/auth/signup/confirm/${1}
}

test_signin(){
    echo "post SignInData"
    curl -v -d "@test/SignInData.json" -H "Content-Type: application/json" POST ${api_url}/auth/signin
}

test_signuptoken(){
    echo "post SignUpToken"
    curl -v -d "@test/SignUpToken.json" -H "Content-Type: application/json" POST ${api_url}/auth/signup/token
}

test_current(){
    echo "get /signup/confirm/:token"
    curl -v GET ${api_url}/v1/signup/confirm/${1}
}


print_help() {
    echo "Commands for drops-backend test"
    echo " full             #full test"
    echo " signup           #signup test"
    echo " confirm          #confirm test"
}
case $1 in 
    full) test_full;;
    signup) test_signup;;
    token) test_signuptoken;;
    signin) test_signin;;
    confirm) test_confirm $2;;
    *) print_help;;
esac
