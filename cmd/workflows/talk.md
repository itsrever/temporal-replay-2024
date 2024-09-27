## Presenting Rever: Launch Fast and Fail Fast
At REVER, established two and a half years ago (YC S22), we manage the return process for e-commerces. While being only 6 backend developers, we have to integrate plenty of third party providers:
- Ecommerce platforms such as Shopify or Woocommerce
- Logistics providers such as UPS or DHL
- Payments providers such as Stripe or Revolut
- ERPs such as SAP or Odoo
- Other ecommerce partners such as 3PLs

In a startup like REVER, being fast in delivering new value is key. In fact, “Launch Fast, Fail Fast” is our main value. This talk is about how Temporal has enabled us to develop a way of working that fits our development speed needs in a very uncertain environment (both businesswise and also technically, as integrations face lots of unexpected challenges). 

## Integrations challenges 
At REVER, when a customer wants to return a product, we offer them the possibility to send someone to pick up the product at their home. To offer this feature, we have to integrate to lots of logistics providers that support pickups in different geographies. But pickups are difficult:
- Providers have heterogeneous constraints and it is sometimes unpredictable to know if a pickup can be scheduled until it actually fails. For example:
  - bank holidays depend on the pickup address
  - no pickups are offered in some time frames (that change constantly)
  - pickups only offered based on current demand
- Some carriers don’t have an api to integrate with and we need to use web automation tools to integrate with them which leaves us vulnerable to changes in the web UI that affect our integration
- Each carrier has different validations for user input (maximum length of the address, for example)

Moreover, anyone who has integrated with third-party services knows that each integration is unique, even within the same industry. APIs are very heterogeneous and can have different levels of maturity. There are some that are stable, well documented, well structured and user friendly, while others change constantly and have undocumented edge cases and hidden behaviours.

Investing in understanding and testing all these behaviours and implement fallbacks for all providers we have to integrate to is in general not justified businesswise for us, as we don’t know how much a given provider will be used, nor how many of these edge cases are actually going to happen, and the cost of opportunity in a startup is just too high.

## ON HOLD: Our way of working
What we we do is just:
1. Integrate the happy path
2. Launch
3. When there is an unexpected failure, the process (workflow) waits for manual action

The part where Temporal comes into play, and the secret sauce for this strategy is step number 3. We call it ON_HOLD and it is treated as a product inside REVER. For the activities we want, we wrap them with a library that calls the activity (with a standard retry policy) and:
- If activity succeeds, returns the output of the activity
- If the activity fails unexpectedly
  - Sets a custom search attribute called ReverStatus to “ON_HOLD”
  - Executes alerting, metrics and auditing activities
  - Waits for a manual action

We then provide our operations team with a dashboard where they can see all the workflows that are ON_HOLD and the error that the activity returned. Then they can act on it depending on the situation:
Potential transient error: they just send a "retry" signal, which is interpreted by the workflow as to run the activity again.
Bug on our side or edge case that we want to solve: deploy a fixed version of the activity worker and send a "retry" signal
Bug on our side or edge case: execute the activity manually (create the pickup in the provider’s UI) and send a "manually-executed" signal, which is interpreted by the workflow as that the activity doesn't need to be executed anymore. The json body of the signal is handled as if it was the output of the activity (for example, the time slot where the pickup will happen)

A simple diagram of this can be seen in this [image](https://i.ibb.co/4jKb5Ck/image-14.png).

To support this with a traditional set up with queues would have been impossible for a team of REVER’s size and it would probably mean having to provide infrastructure for any new activity for which we want to implement this. It would also require a lot of work on the UI side to represent the status or history of each of the flows and also to operate them.

Instead, with Temporal, the implementation inside the workflow is very simple (as can be seen in this [example library](https://github.com/itsrever/temporal-replay-2024) we implemented) and the internal dashboard has been also very easy to develop just using the temporal sdk to filter workflows by a custom search attribute. Also, it is build in such a way so that it requires 0 extra work for our developers to make a new activity to work like this, as everything is activity independent. Out of the box out developers gets all these features, making it a no-brainer to use and making us use it even in very simple activities, removing all fear of launching fast and enhancing experimentation. The little upfront work and these cultural implications for the tech team make it very suitable for other early startups where speed is key. 

A short video demonstration with code samples from the [library](https://github.com/itsrever/temporal-replay-2024) can be found [here](https://www.loom.com/share/66a48dccc21547148a240be96a1f5242).


## But it doesn’t end here
As said, out of the box, we also have metrics on all activities that go to ON_HOLD. Weekly we:
1. Assess together with our operations team which errors give them more work and we estimate the tech effort to solve them 
2. Prioritise the bigger issues in our roadmap
3. Implement the necessary improvements which generally include: improve error handling, manage exception flows, and prevent known error paths

## Some data
This [plot](https://i.ibb.co/YyQ6Ws5/image2.png) shows the percentage of return processes that have at least one ON_HOLD. This includes not only pickup activities, but also all other integrations with payments providers, ecommerce platforms, etc.

We can see a steady decrease in the percentage which is thanks to our team prioritising fixes based on impact. This is especially significant taking into account that REVER is steadily growing 500% year on year.

We can also see increases due to adding new integrations or changes in old ones. 

At the time of writing, this value is around 1.8%. While this still means a lot of manual work for our operations team and is not very scalable at REVER’s actual growth, it is the result of the tradeoff of what we can handle as a team while focusing on delivering value in other areas. In other words, businesswise, we are comfortable with this amount of manual work but surely, in the near future we will need to keep investing in lowering this percentage again. 

As a conclusion, Temporal has enabled us to create this operations framework with a small effort completely assumable for a small team, and this has led REVER to succeed and process returns in around 80 countries with many different partners.
