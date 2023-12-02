import http from "k6/http";
import {check, group} from "k6";
import {
    BASE_URL,
    closeAccount,
    deposit,
    getAccount,
    getAllAccounts,
    getToken,
    requestConfigWithTag,
    withdraw
} from "./helpers.js";

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


let token = ""
let accountId = ""

function createAccount() {
    const res = http.post(`${BASE_URL}/account`,
        JSON.stringify({
            type: "checking",
        }),
        requestConfigWithTag("createAccount", token),
    );

    check(res, {
        "status is 201": (r) => r.status === 201,
        "accountID is present": (r) => r.json("accountID") !== "",
        "response is in json": (r) => r.headers["Content-Type"] === "application/json; charset=utf-8",
    });
    accountId = res.json("accountID");
}

export default () => {
    if (__ITER === 0) {
        group("get token", () => {
            token = getToken();
        });

        group("create account", () => {
            createAccount();
        });
    }

    group("get all accounts", () => {
        getAllAccounts(token);
    });

    group("get account", () => {
        getAccount(accountId, token);
    })

    group("deposit", () => {
        if (!deposit(accountId, token)) {
            return;
        }
    })

    group("withdraw", () => {
        withdraw(accountId, token);
    })

    if (__ITER === ITERATIONS - 1) {
        group("close account", () => {
            closeAccount(accountId, token);
        });
    }
}
