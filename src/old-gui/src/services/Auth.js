const axios = require('axios').default;

// checks if the bot has been authenticated or not
export function authenticated() {
    var authed = false
    axios.get("http://" + window.location.hostname + ":8080/checkauth").then((response) => {
        authed = response.data;
        console.log("in request:", authed)
      }).catch(function(error) {
        console.log(error)
      });

      return authed;
}

// sends a newly generated oauth token to the backend
export function sendToken(hash) {
    var access_token = hash.substr(hash.search(/(?<=^|&)access_token=/))
                    .split('&')[0]
                    .split('=')[1];
  
    axios.post("http://" + window.location.hostname + ":8080/addoauth", JSON.stringify(access_token)).then((res) => {
        console.log(res.data)
      }).catch((error) => {
          console.log(error)
      });
}