import { Client, DefaultSession, TrackingType, InteractiveSession } from "@sajari/sdk-js";
  // new SajariSDK.Client("project", "collection")...
  const websitePipeline = new Client("1627439095818242237", "animal-search").pipeline(
    "website"
  );

  const clickTrackedSession = new InteractiveSession(
    "q",
    new DefaultSession(TrackingType.Click, "url")
  );

  const values = { q: "wings" };
  websitePipeline
    .search(values, clickTrackedSession.next(values))
    .then(([response, values]) => {
      // handle response 
      response.results.forEach((r) => {
        console.log(r)
        const title = document.createElement("a");
        title.textContent = r.values.title;
        title.href = r.values.url;
        title.onmousedown = () => {
          title.href = r.token.click;
        };

        document.body.appendChild(title);
      });
    })
    .catch((error) => {
      // handle error
      console.log(error)
    });

    // most primitive layer between the frontend and backend 
    // hand sdk your account id and pipeline id in wrapper object
    // ask the wrapper to submit a search to the backend
    // customer would have to be an engineer to use it 
    // some customers use vue 
    // there's a pipeline class 
    // react hooks is confusing > 
