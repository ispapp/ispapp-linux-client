import { write } from 'bun';

async function fetchArchs() {
    let page = 1;
    let results: any[] = [];
    let hasMore = true;

    while (hasMore) {
        const response = await fetch(`https://hub.docker.com/v2/repositories/openwrt/sdk/tags?page_size=25&page=${page}&ordering=last_updated&name=23.05`, {
            "credentials": "include",
            "headers": {
                "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:133.0) Gecko/20100101 Firefox/133.0",
                "Accept": "application/json",
                "Accept-Language": "en-US,en;q=0.5",
                "content-type": "application/json",
                "x-docker-api-client": "docker-hub/v5305.0.0",
                "Sec-GPC": "1",
                "Sec-Fetch-Dest": "empty",
                "Sec-Fetch-Mode": "cors",
                "Sec-Fetch-Site": "same-origin",
                "Priority": "u=6",
                "Cookie": "_csrf=IlkxUHBfcktpYzdvLXJZS0NFNEhDci1xTUxiRldPMEppIg==.Qa3EmiRrJssWqEPKh4X6MHJ6eu+Udms39GdZA5c3h9Y; _vis_opt_test_cookie=1; dds-theme={\"preference\":\"system\",\"resolved\":\"dark\"}; fullstoryStart=false; OptanonAlertBoxClosed=2024-10-29T20:49:21.588Z; OptanonConsent=isGpcEnabled=1&datestamp=Tue+Nov+26+2024+07:20:06+GMT+0100+(GMT+01:00)&version=202306.2.0&browserGpcFlag=1&isIABGlobal=false&hosts=&consentId=569774d9-6b51-4bcb-be9b-ffc275c96733&interactionCount=1&landingPath=NotLandingPage&groups=C0003:1,C0001:1,C0002:1,C0004:1&geolocation=;&AwaitingReconsent=false; userty.core.s.fe7522=__WQiOiI3OWI2ODNhODE4MGIwMWY0N2QyMGI3MTEyM2U3MzAyOSIsInN0IjoxNzI1NTI4NjU3NjIwLCJyZWFkeSI6dHJ1ZSwid3MiOiJ7XCJ3XCI6MTUyNCxcImhcIjo3MTB9IiwicHYiOjI0fQ==eyJza"
            },
            "referrer": "https://hub.docker.com/r/openwrt/sdk/tags",
            "method": "GET",
            "mode": "cors"
        });
        if (response.status !== 200) {
            console.error(`Failed to fetch page ${page}: ${response.status} ${response.statusText}`);
            break;
        }
        const data = await response.json();
        results = results.concat(data.results);
        const snapshotResults = data.results.filter((result: { name: string; }) => result.name.includes("23.05"));
        if (snapshotResults.length > 0) {
            console.log("Found results ending with 23.05-SNAPSHOT:", snapshotResults);
        }

        if (data.results.length === 0) {
            hasMore = false;
        } else {
            page++;
        }
    }

    await write('matched_archs.json', JSON.stringify(results, null, 2));
}

fetchArchs();