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
    let username = document.getElementById('username').value
    let expdate = document.getElementById('expdir').value
    expdate = expdate + "T12:00:00.123Z"
    req.open("GET", `http://${website}/setExpireDate?username=${username}&sip=${serverip}&expdate=${expdate}`);
    req.send();
    req.onload = function() {
      if (req.readyState === req.DONE) {
        if (req.status === 200) {
          alert('Changed Expiration Date');
        } else {
          alert('unexpected error')
        }
      }
    };
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
          alert('Successful Login');
        }
      }
    }
  };
}

function showUsers() {
  var userTable = document.getElementById('userTable');
  userTable.innerHTML = '<table class="table table-borderless" id="userTable" style="color: white;"><thead><tr><th scope="col">Username</th><th scope="col">Expiration</th><th scope="col">Last Login</th></tr></thead><tbody></tbody></table>';
  var req = new XMLHttpRequest();
  let website = window.location.href.split('/')[2]
  let serverip = document.getElementById('serverip').value
  if (serverip == undefined || serverip.indexOf(' ')) {
    alert('Please ensure you have set the Server IP, Hub User, and Hub Password properly')
    return
  }
  req.open("GET", `http://${website}/listUsers?sip=${serverip}`);
  req.send();
  req.onload = function() {
    if (req.readyState === req.DONE) {
      if (req.status === 200) {
        let response = JSON.parse(req.response)
        var list = response['result']['UserList']
        list.forEach(user => {
          let username = user['Name_str'];
          let expdate = user['Expires_dt'];
          let lastlogin = user['LastLoginTime_dt'];
          let row = userTable.insertRow(1);
          let cell1 = row.insertCell(0).innerHTML = username;
          let cell2 = row.insertCell(1).innerHTML = expdate.replace('T', ' ').replace('Z', '');
          let cell3 = row.insertCell(2).innerHTML = lastlogin.replace('T', ' ').replace('Z', '');
        });
      }
    }
  }
}
