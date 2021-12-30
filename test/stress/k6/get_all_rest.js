import http from 'k6/http';
import { check } from 'k6';
import { Rate } from 'k6/metrics';

export const errorRate = new Rate('errors');

const accessToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDA4MzIyMjgsInVzZXJuYW1lIjoibmFqaWUifQ.kxe_VBTaE7LAhO8YlzwIVCOCv0weLK0Hyu2VvtqGXs8';

export default function () {
    const url = 'http://localhost:8090/links';
    const params = {
        headers: {
            'Authorization': `Bearer ${accessToken}`,
            'Content-Type': 'application/json',
        },
    };
    check(http.get(url, params), {
        'status is 200': (r) => r.status == 200,
    }) || errorRate.add(1);
}