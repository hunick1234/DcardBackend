const url = "http://localhost:8080/api/v1/ad?";

export const generateGetAd = () => {
  return {
    method: "GET",
    url: url + generateGetQuery(),
  };
};

export const generatePostAd = (i) => {
  return {
    method: "POST",
    url: url,
    payload: generateAdPayload(i),
    headers: { "Content-Type": "application/json" },
  };
};
const countries = ["US", "JP", "TW", "KR", "CN", "SG", "AD", "AE"];
const genders = ["M", "F"];
const platforms = ["android", "ios", "web"];

function generateAdPayload(i) {
  const startDate = new Date("2022-01-01");
  const endDate = new Date("2026-12-31");
  const startAt = randomDate(startDate, endDate);
  const endAt = new Date(
    startAt.getTime() + randomNumber(1, 5) * 24 * 60 * 60 * 1000
  ); // Ensure endAt is after startAt

  return JSON.stringify({
    title: "ad " + `${i}`,
    startAt: startAt.toISOString(),
    endAt: endAt.toISOString(),
    conditions: {
      ageStart: randomNumber(1, 50),
      ageEnd: randomNumber(50, 100),
      country: randomItem(countries),
      gender: randomItem(genders),
      platform: randomItem(platforms),
    },
  });
}

function randomItem(arr) {
  //rendom between 0 and arr.lenth\
  let result = [];
  let rand = Math.floor(Math.random() * arr.length);
  for (let i = 0; i < rand; i++) {
    result.push(arr[Math.floor(Math.random() * arr.length)]);
  }
  return result;
}

function randomNumber(min, max) {
  return Math.floor(Math.random() * (max - min + 1) + min);
}

function randomDate(start, end) {
  return new Date(
    start.getTime() + Math.random() * (end.getTime() - start.getTime())
  );
}

function generateGetQuery() {
  // Generating random values
  let uri = "";
  const offset = randomNumber(0, 100); // Assuming offset range
  uri += `offset=${offset}&`;
  const limit = randomNumber(1, 100); // Assuming limit range
  uri += `limit=${limit}&`;
  let country = randomItem(countries);
  if (country.length !== 0) {
    country = country.join(",");
    uri += `country=${country}&`;
  }
  let gender = randomItem(genders);
  if (gender.length !== 0) {
    gender = gender.join(",");
    uri += `gender=${gender}&`;
  }

  let platform = randomItem(platforms);
  if (platform.length !== 0) {
    platform = platform.join(",");
    uri += `platform=${platform}&`;
  }
  // Constructing query string
  return uri;
}
