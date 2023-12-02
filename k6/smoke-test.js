import http from "k6/http";
import {check, group} from "k6";

const ITERATIONS = 10;

export const options = {
    scenarios: {
        smoke_test: {
            executor: "per-vu-iterations",
            vus: 3,
            iterations: ITERATIONS,
        },
    },
};


const BASE_URL = "http://localhost:8080/api/v1";
let token = ""
let accountId = ""

const requestConfigWithTag = (tag) => ({
    headers: {
        Content_Type: "application/json",
        Authorization: `Bearer ${token}`,
    },

    tags: Object.assign(
        {name: tag},
    ),
});

function getToken() {
    const res = http.get(`${BASE_URL}/login`);
    check(res, {
        "status is 200": (r) => r.status === 200,
        "token is present": (r) => r.json("token") !== "",
        "response is in json": (r) => r.headers["Content-Type"] === "application/json; charset=utf-8",
    });
    token = res.json("token");
}

function createAccount() {
    const res = http.post(`${BASE_URL}/account`,
        JSON.stringify({
            type: "checking",
        }),
        requestConfigWithTag("createAccount"),
    );

    check(res, {
        "status is 201": (r) => r.status === 201,
        "accountID is present": (r) => r.json("accountID") !== "",
        "response is in json": (r) => r.headers["Content-Type"] === "application/json; charset=utf-8",
    });
    accountId = res.json("accountID");
}

function getAllAccounts() {
    const res = http.get(`${BASE_URL}/accounts/all`, requestConfigWithTag("getAllAccounts"));

    check(res, {
        "status is 200": (r) => r.status === 200,
        "response is in json": (r) => r.headers["Content-Type"] === "application/json; charset=utf-8",
    });
}

function getAccount() {
    const res = http.get(`${BASE_URL}/account/${accountId}`, requestConfigWithTag("getAccount"));

    check(res, {
        "status is 200": (r) => r.status === 200,
        "response is in json": (r) => r.headers["Content-Type"] === "application/json; charset=utf-8",
    });
}

function randomFloat(min, max) {
    return +(Math.random() * (max - min + 1) + min).toFixed(2);
}

function deposit() {
    const res = http.patch(`${BASE_URL}/account/${accountId}/deposit`,
        JSON.stringify({
            amount: randomFloat(90, 200),
        }),
        requestConfigWithTag("deposit"),
    );

    const ok = check(res, {
        "status is 204": (r) => r.status === 204,
    });

    if (!ok) {
        console.log(`deposit failed: ${res.status} ${res.body}`);
    }
    return ok;
}

function withdraw() {
    const res = http.patch(`${BASE_URL}/account/${accountId}/withdraw`,
        JSON.stringify({
            amount: randomFloat(20, 80),
        }),
        requestConfigWithTag("withdraw"),
    );

    const ok = check(res, {
        "status is 204": (r) => r.status === 204,
    });

    if (!ok) {
        console.log(`withdraw failed: ${res.status} ${res.body}`);
    }
    return ok;
}

function closeAccount() {
    const res = http.patch(`${BASE_URL}/account/${accountId}/close`, requestConfigWithTag("closeAccount"),
        requestConfigWithTag("withdraw"));

    const ok = check(res, {
        "status is 204": (r) => r.status === 204,
    });

    if (!ok) {
        console.log(`closeAccount failed: ${res.status} ${res.body}`);
    }
    return ok;
}

export default () => {
    if (__ITER === 0) {
        group("get token", () => {
            getToken();
        });

        group("create account", () => {
            createAccount();
        });
    }

    group("get all accounts", () => {
        getAllAccounts();
    });

    group("get account", () => {
        getAccount();
    })

    group("deposit", () => {
        if (!deposit()) {
            return;
        }
    })

    group("withdraw", () => {
        withdraw();
    })

    if (__ITER === ITERATIONS - 1) {
        group("close account", () => {
            closeAccount();
        });
    }
}
