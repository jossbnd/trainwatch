import Toybox.Graphics;
import Toybox.Lang;
import Toybox.WatchUi;

class GarminView extends WatchUi.View {

    private var _train1 as Lang.String = "";
    private var _train2 as Lang.String = "";
    private var _status as Lang.String = "Press Start";

    function initialize() {
        View.initialize();
    }

    function onLayout(dc as Dc) as Void {
        setLayout(Rez.Layouts.MainLayout(dc));
    }

    function onShow() as Void {
    }

    function onUpdate(dc as Dc) as Void {
        (findDrawableById("LabelTrain1") as WatchUi.Text).setText(_train1);
        (findDrawableById("LabelTrain2") as WatchUi.Text).setText(_train2);
        (findDrawableById("LabelStatus") as WatchUi.Text).setText(_status);
        View.onUpdate(dc);
    }

    function onHide() as Void {
    }

    function SetLoading() as Void {
        _train1 = "";
        _train2 = "";
        _status = "Loading...";
        WatchUi.requestUpdate();
    }

    function SetError(message as Lang.String) as Void {
        _train1 = "";
        _train2 = "";
        _status = message;
        WatchUi.requestUpdate();
    }

    function SetTrains(train1 as Lang.String, train2 as Lang.String) as Void {
        _train1 = train1;
        _train2 = train2;
        _status = "";
        WatchUi.requestUpdate();
    }
}
