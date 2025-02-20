"use strict";

function indicateError(message, sec=30000) {
  document.getElementById("error-detail").innerHTML = message;
  document.getElementById("error-container").className = "error";

  setTimeout(function() {
    document.getElementById("error-container").className = "error non-display";
  }, sec);
}

function indicateCheck(message, sec=3000) {
    document.getElementById("check-detail").innerHTML = message;
    document.getElementById("check-container").className = "check";

    setTimeout(function() {
      document.getElementById("check-container").className = "check non-display";
    }, sec);
}

async function checkLogin() {
  const url = "http://localhost:8080/status/Islogin";
  const response = await fetch(url);
  if (!response.ok) {
    indicateError(`HTTP Error: ${response.status}`);
    return false;
  } else {
    const state = await response.json();
    if (!state.IsLogin) {
      indicateError(`<a href="./login">ログイン</a>していません`);
    }
    return state.isLogin;
  }
}

async function getPlayingState() {
  const url = "http://localhost:8080/status/playing";
  try {
    const response = await fetch(url);
    if (!response.ok) {
      indicateError(`HTTP Error: ${response.status}`);
    }
    const state = await response.json();
    console.log(state.message);
    if (state.message == "there is no playing state") {
      document.getElementById("artwork").src = state.artwork;
      document.getElementById("title").innerHTML = "再生停止中";
      document.getElementById("artist").innerHTML = `&ndash;`;
      document.getElementById("album").innerHTML = `&ndash;`;
      document.getElementById("link").href = `https://spotify.com/`;
    } else {
      document.getElementById("artwork").src = state.Image;
      document.getElementById("title").innerHTML = state.Name;
      document.getElementById("artist").innerHTML = state.ArtistsPureString;
      document.getElementById("album").innerHTML = state.Album;
      document.getElementById("link").href = `https://open.spotify.com/intl-ja/track/${state.ID}`;
    }
  } catch (error) {
    indicateError(`Error: ${error.message}`);
  }
  return 0;
}

async function getPositionState() {
  if (!navigator.geolocation) {
    indicateError(`位置情報サービスはあなたの端末でサポートされていません`);
    return 0;
  } else {
    console.log("Locating…");
    navigator.geolocation.getCurrentPosition(
      async function (position) {
        document.getElementById("lon").innerHTML = position.coords.longitude.toPrecision(6);
        document.getElementById("lat").innerHTML = position.coords.latitude.toPrecision(6);
        const url = `http://localhost:8080/get/location?lat=${position.coords.latitude}&lon=${position.coords.longitude}`;
        const response = await fetch(url);
        try {
          if (!response.ok) {
            indicateError(`HTTP Error: ${response.status}`);
          }
          const positionList = await response.json();
          for(let i=0;i<positionList.length;i++) {
            document.getElementById(`place-${i+1}`).innerHTML = positionList[i].Name;
            document.getElementById(`area-${i+1}`).innerHTML = positionList[i].AreaName;
            document.getElementById(`coords-${i+1}`).innerHTML = `緯度: ${positionList[i].Latitude.toPrecision(6)}, 経度: ${positionList[i].Longitude.toPrecision(6)}`;
          }
        } catch (error) {
          indicateError(`Error: ${error.message}`);
        }
      },
      function () {
        indicateError(`位置情報を取得できませんでした`);
      }
    )
  }
}


let res = checkLogin()
if (res == true) {
  getPlayingState();
  getPositionState();
};
