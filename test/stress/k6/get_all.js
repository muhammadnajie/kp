import http from "k6/http";
import { sleep } from "k6";

export default function() {
    const accessToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDA4MzIyMjgsInVzZXJuYW1lIjoibmFqaWUifQ.kxe_VBTaE7LAhO8YlzwIVCOCv0weLK0Hyu2VvtqGXs8';

    const query = `
    query {
        links {
            title
            address
        }
    }`;

    const headers = {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
    };

    const res = http.post('http://localhost:8090/query', JSON.stringify({ query: query }), {
        headers: headers,
    });

    // if (res.status === 200) {
    //     console.log(JSON.stringify(res.body));
    //     const body = JSON.parse(res.body);
    //     const issue = body.data.links;
    //     console.log(issue);
    // }
}

