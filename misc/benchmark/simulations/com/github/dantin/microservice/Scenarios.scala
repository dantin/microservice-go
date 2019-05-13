package com.github.dantin.microservice

import scala.concurrent.duration._

import io.gatling.core.Predef._
import io.gatling.http.Predef._
import io.gatling.jdbc.Predef._

object Scenarios {

  val rampUpTimeSecs = 10

  var browse_guids = csv("accounts.csv").circular

  val scn_Browser = scenario("GetAccounts")
    .during(Conf.duration) {
      feed(browse_guids)
        .exec(
          http("GetAccount")
            .get("/accounts/" + "${accountId}")
            .headers(Headers.http_header)
            .check(status.is(200))
          ).pause(1)
    }
}
