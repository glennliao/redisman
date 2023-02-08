export function request(url: string, method: string, data: Record<string, any>) {
  url = `http://${window.basicURL}${url}`;
  return fetch(url, { credentials: "include", method, body: JSON.stringify(data) })
    .then(resp => resp.json())
}
