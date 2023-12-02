import http from "k6/http";
import {check} from "k6";

export const BASE_URL = "http://localhost:8080/api/v1";

export const requestConfigWithTag = (tag, token) => ({
    headers: {
        Content_Type: "application/json",
        Authorization: `Bearer ${token}`,
    },

    tags: Object.assign(
        {name: tag},
    ),
});

export function getToken() {
    const res = http.get(`${BASE_URL}/login`);
    const ok = check(res, {
        "status is 200": (r) => r.status === 200,
        "token is present": (r) => r.json("token") !== "",
        "response is in json": (r) => r.headers["Content-Type"] === "application/json; charset=utf-8",
    });

    if (!ok) {
        console.log(`token request failed: ${res.status} ${res.body}`);
        return "";
    }
    return res.json("token");
}

export function getAllAccounts(token) {
    const res = http.get(`${BASE_URL}/accounts/all`, requestConfigWithTag("getAllAccounts", token));

    const ok = check(res, {
        "status is 200": (r) => r.status === 200,
        "response is in json": (r) => r.headers["Content-Type"] === "application/json; charset=utf-8",
    });

    if (!ok) {
        console.log(`get all accounts failed: ${res.status} ${res.body}`);
    }
}

export function getAccount(accountId, token) {
    const res = http.get(`${BASE_URL}/account/${accountId}`, requestConfigWithTag("getAccount", token));

    const ok = check(res, {
        "status is 200": (r) => r.status === 200,
        "response is in json": (r) => r.headers["Content-Type"] === "application/json; charset=utf-8",
    });

    if (!ok) {
        console.log(`get account failed: ${res.status} ${res.body}`);
    }
}

function randomFloat(min, max) {
    return +(Math.random() * (max - min + 1) + min).toFixed(2);
}

export function deposit(accountId, token) {
    const res = http.patch(`${BASE_URL}/account/${accountId}/deposit`,
        JSON.stringify({
            amount: randomFloat(90, 200),
        }),
        requestConfigWithTag("deposit", token),
    );

    const ok = check(res, {
        "status is 204": (r) => r.status === 204,
    });

    if (!ok) {
        console.log(`deposit failed: ${res.status} ${res.body}`);
    }
    return ok;
}

export function withdraw(accountId, token) {
    const res = http.patch(`${BASE_URL}/account/${accountId}/withdraw`,
        JSON.stringify({
            amount: randomFloat(20, 80),
        }),
        requestConfigWithTag("withdraw", token),
    );

    const ok = check(res, {
        "status is 204": (r) => r.status === 204,
    });

    if (!ok) {
        console.log(`withdraw failed: ${res.status} ${res.body}`);
    }
}

export function closeAccount(accountId, token) {
    const res = http.patch(`${BASE_URL}/account/${accountId}/close`, null, requestConfigWithTag("closeAccount", token));

    const ok = check(res, {
        "status is 204": (r) => r.status === 204,
    });

    if (!ok) {
        console.log(`close account failed: ${res.status} ${res.body}`);
    }
}

export function deleteAccount(accountId, token) {
    const res = http.del(`${BASE_URL}/account/${accountId}`, null, requestConfigWithTag("deleteAccount", token));

    const ok = check(res, {
        "status is 204": (r) => r.status === 204,
    });

    if (!ok) {
        console.log(`delete account failed: ${res.status} ${res.body}`);
    }
}