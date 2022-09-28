import http from "k6/http";
import { check, group, sleep } from "k6";
import { Rate } from "k6/metrics";

export const baseUrl = "http://127.0.0.1:3000";

function createUser(email, pass) {
  const url = baseUrl + "/v1/auth/register";
  const params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  const body = JSON.stringify({
    email: email,
    password: pass,
  });

  check(http.post(url, body, params), {
    "no duplicated credentials": (r) => r.status != 409,
    "created new user": (r) => r.status <= 299,
  });
}

function login(email, pass) {
  const url = baseUrl + "/v1/nc/auth/login";
  const params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  const body = JSON.stringify({
    email: email,
    password: pass,
  });

  var userID = 0;
  var sessionID = "";
  let resp = http.post(url, body, params);

  check(resp, {
    "logged in": (r) => r.status <= 299,
  });

  if (resp.status <= 299) {
    userID = resp.json("data.id");
    sessionID = resp.headers["Session-Id"];
  }

  return [userID, sessionID];
}

function deleteUser(id, sessionID) {
  const url = baseUrl + "/v1/nc/users/" + id;
  const params = {
    headers: {
      "Content-Type": "application/json",
      "Session-Id": sessionID,
    },
  };

  let resp = http.del(url, {}, params)
  check(resp, {
    "user deleted": (r) => r.status <= 299,
  });

  
  if (resp.status > 299) {
    console.log(resp)
  }
}

function getAllProduct(sessionID) {
  const url = baseUrl + "/v1/nc/products";
  const params = {
    headers: {
      "Content-Type": "application/json",
      "Session-Id": sessionID,
    },
  };

  let resp = http.get(url, params);
  check(resp, {
    "product fetched": (r) => r.status <= 299,
  });

  if (resp.status <= 299) {
    return resp.json("data");
  } else {
    console.log(resp);
  }
}

function getProductByID(pid, sesssionID) {
  const url = baseUrl + "/v1/nc/products/" + pid;
  const params = {
    headers: {
      "Content-Type": "application/json",
      "Session-Id": sesssionID,
    },
  };

  let resp = http.get(url, params);

  check(resp, {
    "product fetched": (r) => r.status <= 299,
  });

  if (resp.status <= 299) {
    return resp.json("data");
  }
}

function checkout(data, sessionID) {
  const url = baseUrl + "/v1/nc/checkout";
  const params = {
    headers: {
      "Content-Type": "application/json",
      "Session-Id": sessionID,
    },
  };

  let resp = http.post(url, JSON.stringify({ items: data }), params);

  check(resp, {
    "product checkout'd": (r) => r.status <= 299,
  });

  if (resp.status > 299) {
    console.log(resp)
  }
}

export default function () {
  let userID = 0;
  let sessionID = "";
  let email = __VU + __ITER + Math.floor(Math.random() * 1000) + "@mairu.xyz";
  let password = "1111";
  let itemCart = [];

  group("initialize user account", () => {
    createUser(email, password);
    [userID, sessionID] = login(email, password);
  });
  sleep(0.3);
  group("browsing marketplace", () => {
    let products = getAllProduct(sessionID);
    let id = products[Math.floor(Math.random() * products.length)].product_id;
    let product = getProductByID(id, sessionID);
    itemCart.push({
      product_id: product.product_id,
      unit_price: product.unit_price,
      qty: Math.floor(Math.random() * 10) + 1,
    });
  });
  sleep(0.4);
  group("checkout product", () => {
    [userID, sessionID] = login(email, password);
    sleep(0.3);
    console.log(userID, sessionID)
    checkout(itemCart, sessionID);
  });
  sleep(0.4);
  group("post cleanup", () => {
    [userID, sessionID] = login(email, password);
    sleep(0.3);
    if (userID != 0) {
      deleteUser(userID, sessionID);
    } else {
      console.log(userID, sessionID, email)
    }
  })

}
