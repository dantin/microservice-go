package com.github.dantin.microservice

import scala.concurrent.duration._

import io.gatling.core.Predef._
import io.gatling.http.Predef._
import io.gatling.jdbc.Predef._

class LoadTest extends Simulation {
  setUp(
    Scenarios.scn_Browser.inject(
      rampUsers(Conf.users) duration (10 seconds)).protocols(
        Conf.httpConf)
    )
}
