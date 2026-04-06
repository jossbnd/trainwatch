import Toybox.Application;
import Toybox.Lang;
import Toybox.WatchUi;

class GarminApp extends Application.AppBase {

    function initialize() {
        AppBase.initialize();
    }

    function onStart(state as Dictionary?) as Void {
    }

    function onStop(state as Dictionary?) as Void {
    }

    function getInitialView() as [Views] or [Views, InputDelegates] {
        var view = new GarminView();
        var delegate = new GarminDelegate(view);
        return [view, delegate];
    }

}

function getApp() as GarminApp {
    return Application.getApp() as GarminApp;
}
