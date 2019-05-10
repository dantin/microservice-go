package com.github.dantin.microservice

import io.gatling.core.Predef._
import io.gatling.http.Predef._
import io.gatling.jdbc.Predef._

object Conf {
  var users = System.getProperty("users", "2").toInt
  var baseUrl = System.getProperty("baseUrl", "http://192.168.0.119:30000")
  var httpConf = http.baseUrl(baseUrl)
  var duration = System.getProperty("duration", "30").toInt
}
