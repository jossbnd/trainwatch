import Toybox.Lang;
import Toybox.WatchUi;
import Toybox.Communications;
import Toybox.Time;
import Toybox.Time.Gregorian;
import Toybox.System;

class GarminDelegate extends WatchUi.BehaviorDelegate {

    private const API_URL  = Config.API_URL;
    private const API_KEY  = Config.API_KEY;
    private const STOP_REF = Config.STOP_REF;
    private const LINE_REF = Config.LINE_REF;

    private var _view as GarminView;
    private var _loading as Lang.Boolean = false;
    private var _lastFetchMs as Lang.Number = -5000;

    function initialize(view as GarminView) {
        BehaviorDelegate.initialize();
        _view = view;
    }

    function onSelect() as Lang.Boolean {
        var now = System.getTimer();
        if (_loading || (now - _lastFetchMs) < 2000) { return true; }
        FetchDepartures();
        return true;
    }

    private function FetchDepartures() as Void {
        _lastFetchMs = System.getTimer();
        _loading = true;
        _view.SetLoading();

        var params = {
            "stop_ref" => STOP_REF,
            "line_ref" => LINE_REF,
            "limit"    => 2
        };

        var options = {
            :method       => Communications.HTTP_REQUEST_METHOD_GET,
            :headers      => {
                "X-API-Key" => API_KEY,
                "Accept"    => "application/json"
            },
            :responseType => Communications.HTTP_RESPONSE_CONTENT_TYPE_JSON
        };

        Communications.makeWebRequest(API_URL, params, options, method(:OnResponse));
    }

    function OnResponse(responseCode as Lang.Number, data as Lang.Dictionary?) as Void {
        if (responseCode != 200 || data == null) {
            _loading = false;
            _view.SetError("Error " + responseCode);
            return;
        }

        var departures = data["departures"] as Lang.Array;
        if (departures.size() == 0) {
            _loading = false;
            _view.SetError("No departures");
            return;
        }

        var dep1 = FormatDeparture(departures[0] as Lang.Dictionary);
        var dep2 = departures.size() > 1 ? FormatDeparture(departures[1] as Lang.Dictionary) : "";
        _loading = false;
        _view.SetTrains(dep1, dep2);
    }

    private function FormatDeparture(departure as Lang.Dictionary) as Lang.String {
        var destination = departure["destination"] as Lang.String;
        var estimatedAt = departure["estimated_at"] as Lang.String;
        var delayMin    = departure["delay_minutes"] as Lang.Number;

        var minutes = MinutesFromNow(estimatedAt);
        var label   = destination + " " + minutes + "min";

        if (delayMin > 0) {
            label = label + " +" + delayMin + "min";
        }

        return label;
    }

    private function MinutesFromNow(isoTime as Lang.String) as Lang.Number {
        var now   = Time.now().value();
        var year  = (isoTime.substring(0, 4)   as Lang.String).toNumber() as Lang.Number;
        var month = (isoTime.substring(5, 7)   as Lang.String).toNumber() as Lang.Number;
        var day   = (isoTime.substring(8, 10)  as Lang.String).toNumber() as Lang.Number;
        var hour  = (isoTime.substring(11, 13) as Lang.String).toNumber() as Lang.Number;
        var min   = (isoTime.substring(14, 16) as Lang.String).toNumber() as Lang.Number;
        var sec   = (isoTime.substring(17, 19) as Lang.String).toNumber() as Lang.Number;

        var moment = Gregorian.moment({
            :year   => year,
            :month  => month,
            :day    => day,
            :hour   => hour,
            :minute => min,
            :second => sec
        });

        var diff = moment.value() - now;
        if (diff < 0) { diff = 0; }
        return (diff / 60).toNumber();
    }
}
