import http from "k6/http";
import {check, group} from "k6";
import {
    BASE_URL,
    closeAccount,
    deleteAccount,
    deposit,
    getAccount,
    getAllAccounts,
    getToken,
    requestConfigWithTag,
    withdraw
} from "./helpers.js"

export const options = {
    scenarios: {
        stress_test: {
            executor: "ramping-arrival-rate",
            stages: [
                {target: 100, duration: "10s"},
                {target: 100, duration: "20s"},
                {target: 500, duration: "1m"},
                {target: 500, duration: "40s"},
                {target: 0, duration: "20s"},
            ],
            preAllocatedVUs: 230,
        },
    },
};

function createAccounts(token) {
    const req1 = {
        method: "POST",
        url: `${BASE_URL}/account`,
        body: JSON.stringify({
            type: "checking",
        }),
        params: requestConfigWithTag("createAccount", token)
    };

    const req2 = {
        method: "POST",
        url: `${BASE_URL}/account`,
        body: JSON.stringify({
            type: "saving",
        }),
        params: requestConfigWithTag("createAccount", token)
    };

    const responses = http.batch([req1, req2]);

    for (let res of responses) {
        const ok = check(res, {
            "status is 201": (r) => r.status === 201,
            "accountID is present": (r) => r.json("accountID") !== "",
            "response is in json": (r) => r.headers["Content-Type"] === "application/json; charset=utf-8",
        });

        if (!ok) {
            console.log(`create account failed: ${res.status} ${res.body}`);
            return ["", ""];
        }
    }

    return [responses[0].json("accountID"), responses[1].json("accountID")];
}

export default () => {
    const token = group("01. Get token", () => {
        return getToken();
    });

    let [id1, id2] = group("02. Create accounts", () => {
        let [id1, id2] = createAccounts(token);
        if (id1 === "" || id2 === "") {
            return;
        }
        return [id1, id2];
    });

    group("03. Get all accounts", () => {
        getAllAccounts(token);
    });

    group("04. Get accounts", () => {
        getAccount(id1, token);
        getAccount(id2, token);
    })

    group("05. Deposit", () => {
        deposit(id1, token);
        deposit(id2, token);
    })

    group("06. Withdraw", () => {
        withdraw(id1, token);
        withdraw(id2, token);
    })

    group("07. Close account", () => {
        closeAccount(id1, token);
        closeAccount(id2, token);
    });

    group("08. Delete account", () => {
        deleteAccount(id1, token);
        deleteAccount(id2, token);
    });
}
