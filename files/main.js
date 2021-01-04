function AlertInit() {
  let website = window.location.href.split('/')[2]
  let serverip = document.getElementById('serverip').value
  let option = document.getElementById('opt-select').value
  var req = new XMLHttpRequest();
  if (option == 'Create User') {
    let username = document.getElementById('username').value
    let password = btoa(document.getElementById('password').value)
    req.open("GET", `http://${website}/createUser?username=${username}&password=${password}&sip=${serverip}`);
    req.send();
    req.onload = function() {
      if (req.readyState === req.DONE) {
        if (req.status === 200) {
          var response = JSON.parse(req.response);
          if (response['error']) {
            alert('User already exists')
          } else {
            alert('Created User');
          }
        }
      }
    };
  } else if (option == 'Delete User') {
    let username = document.getElementById('username').value
    req.open("GET", `http://${website}/deleteUser?username=${username}&sip=${serverip}`);
    req.send();
    req.onload = function() {
      if (req.readyState === req.DONE) {
        if (req.status === 200) {
          var response = JSON.parse(req.response);
          if (response['error']) {
            alert('User already deleted or does not exist')
          } else {
            alert('Deleted User');
          }
        }
      }
    };
  } else if (option == 'Change Password') {
    let username = document.getElementById('username').value
    let password = btoa(document.getElementById('password').value)
    req.open("GET", `http://${website}/changePassword?username=${username}&password=${password}&sip=${serverip}`);
    req.send();
    req.onload = function() {
      if (req.readyState === req.DONE) {
        if (req.status === 200) {
          alert("Changed user password")
        }
      }
    }
  } else if (option == 'Change Expiration Date') {
    // Not Yet
  } else if (option == 'View User Info') {
    let username = document.getElementById('username').value
    req.open("GET", `http://${website}/getUser?username=${username}&sip=${serverip}`);
    req.send();
    req.onload = function() {
      if (req.readyState === req.DONE) {
        if (req.status === 200) {
          response = JSON.parse(req.response)
          var userObject = response
          console.log(userObject)
          document.getElementById('results').innerHTML = 'Check the developer\'s console (Ctrl+Shift+I or CMD+Shift+I)'
        }
      }
    };
  } else {
    alert('Please input an option from the Drop Down box')
  }
}

function setAuthKey() {
  let website = window.location.href.split('/')[2]
  var req = new XMLHttpRequest();
  let hubUsername = document.getElementById('hub-username').value
  let hubPassword = document.getElementById('hub-password').value
  req.open("GET", `http://${website}/init?hubuser=${hubUsername}&hubpass=${hubPassword}`);
  req.send();
  req.onload = function() {
    if (req.readyState === req.DONE) {
      if (req.status === 200) {
        var response = req.response
        if (response == 'error') {
          alert('Error')
        } else {
          alert('Set Hub Username and Password');
        }
      }
    }
  };
}
