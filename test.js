function positions() {
    // Parse the JSON text into an object
    const jsonObject = tc.response.json;

    // Create the positions array using the values from the JSON object
    const twostring = jsonObject.positions.map((value) => value);

    // Create the pwd variable
    const onestring = jsonObject.positions.map(Number).join("");

    tc.setVar("onestring", onestring);
    tc.setVar("twostring", twostring);
}
module.exports = [positions];
