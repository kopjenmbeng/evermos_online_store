# evermos-online-store

This api services is answer of problem when 12.12 event, which is
1. Customers were able to put items in their cart, check out, and then pay. After several days, many of our customers received calls from
2. our Customer Service department stating that their orders have been canceled due to stock unavailability.
These bad reviews generally come within a week after our 12.12 event, in which we held a large flash sale and set up other major
discounts to promote our store.

After checking in with our Customer Service and Order Processing departments, we received the following additional facts:
1. Our inventory quantities are often misreported, and some items even go as far as having a negative inventory quantity.
The misreported items are those that performed very well on our 12.12 event.
2. Because of these misreported inventory quantities, the Order Processing department was unable to fulfill a lot of orders, and thus
requested help from our Customer Service department to call our customers and notify them that we have had to cancel their orders.

And here I am trying to find the root of the problem why this is happening and this is the posibility.
- the possibility that occurs is when the system processes purchases, the system does not validate the stock first. so that it can cause an imbalance between stock and demand

So I wil give solution start from add to chart until process order.

1. For add to chart process This is the flow that I recommend.

![Optional Text](../master/document/add_to_chart.jpg)

2. For update to chart process This is the flow that I recommend.

![Optional Text](../master/document/update_chart.jpg)

3. For order process This is the flow that I recommend.

![Optional Text](../master/document/order_process.jpg)

# From the solution that I offered above, the database diagram that I will make is like this.

![Optional Text](../master/document/db_schema.jpg)


# And now we are going to technical section.

# Technolgy Stack
-   GO
-   Postgres
-   Consul for store configuration (Optional)
-   Postman for API Documentation

# Before To RUN
![Please restore this database backup](../master/script/onlinestor_db.backup)