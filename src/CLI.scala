import java.time.{Instant, LocalDate, LocalDateTime, ZoneId, ZonedDateTime}
import java.time.format.DateTimeFormatter

import play.api.libs.json._

import scala.io.Source

class CLI(game: String, team: String = "" , daysAhead: Int = 0) {
  //Makes API request to Pandascore
  def get(url: String): String = scala.io.Source.fromURL(url).mkString
  val gameAPIString: String = "&filter[videogame]=" + game
  val teamAPIString: String = "&search[name]=" + team
  var timeAPIString: String = ""
  //Parses API token from config file
  val tokenParse = Json.parse(Source.fromFile("config.json").mkString)
  val token: String = "&token=" + (tokenParse \ "token").as[String]
  var jsonImport: String = ""

  if(daysAhead > 0) {
    val now = ZonedDateTime.now(ZoneId.systemDefault())
    timeAPIString = "&range[begin_at]=" + now + "," + ZonedDateTime.now(ZoneId.systemDefault()).plusDays(daysAhead)
    jsonImport = get("https://api.pandascore.co/matches/upcoming?" + teamAPIString + gameAPIString + timeAPIString + token)
  }
  else {
    jsonImport = get("https://api.pandascore.co/matches/upcoming?" + teamAPIString + gameAPIString + token)
  }

  //Processes JSON String
  val myFormat = DateTimeFormatter.ofPattern("MM/dd HH:mm")
  val parse: JsValue = Json.parse(jsonImport)

  //Extracts team matchup and event time
  var matchup = (parse \ 0 \ "name").as[String]
  matchup = matchup.substring(matchup.indexOf(':') + 2, matchup.length)
  var matchTime = (parse \ 0 \"begin_at").as[String]
  var matchTimeParsed = ZonedDateTime.parse(matchTime)
  matchTimeParsed = matchTimeParsed.withZoneSameInstant(ZoneId.of("America/New_York"))
  val matchTimeString = matchTimeParsed.format(myFormat)

  println("Matchup: " + matchup + '\n'+
    "Time: " + matchTimeString + '\n')
}

object CLI {
  def main(args: Array[String]): Unit = {
    new CLI(args(0), args(1), args(2).toInt)
  }
}