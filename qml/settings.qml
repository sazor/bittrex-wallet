import QtQuick 2.7
import QtQuick.Controls 2.3
import QtQuick.Layouts 1.3
import QtQuick.Controls.Universal 2.0
import QtQuick.Window 2.3

Page {
    width: 350
    height: 380
    visible: true

    StackLayout {
        id: pages
        anchors.fill: parent
        currentIndex: tabBar.currentIndex

        Pane {
            width: parent.width
            height: parent.height
            Column {
                width: parent.width
                spacing: 20
                Row {
                    width: parent.width

                    Column {
                        width: parent.width * 0.2
                        Label {
                            width: parent.width
                            font.pixelSize: 16
                            verticalAlignment: Text.AlignVCenter
                            lineHeightMode: Text.FixedHeight
                            lineHeight: key.height
                            text: "Key"
                            leftPadding: 5
                        }
                    }

                    Column {
                        width: parent.width * 0.8
                        TextField {
                            id: key
                            width: parent.width
                        }
                    }
                }

                Row {
                    width: parent.width

                    Column {
                        width: parent.width * 0.2
                        Label {
                            width: parent.width
                            font.pixelSize: 16
                            verticalAlignment: Text.AlignVCenter
                            lineHeightMode: Text.FixedHeight
                            lineHeight: secret.height
                            text: "Secret"
                            leftPadding: 5
                        }
                    }

                    Column {
                        width: parent.width * 0.8
                        TextField {
                            id: secret
                            width: parent.width
                        }
                    }
                }

                Row {
                    width: parent.width
                    Column {
                        width: parent.width
                        RoundButton {
                            text: "\u2713"
                            anchors.right: parent.right
                            anchors.rightMargin: 0
                            onClicked: textArea.readOnly = true
                        }
                    }
                }
            }
        }

        Pane {
            width: parent.width
            height: parent.height

            Column {
                spacing: 40
                width: parent.width

                Label {
                    width: parent.width
                    wrapMode: Label.Wrap
                    horizontalAlignment: Qt.AlignHCenter
                    text: "TabBar is a bar with icons or text which allows the user"
                          + "to switch between different subtasks, views, or modes."
                }
            }
        }

        Pane {
            width: parent.width
            height: parent.height

            Column {
                spacing: 40
                width: parent.width

                Label {
                    width: parent.width
                    wrapMode: Label.Wrap
                    horizontalAlignment: Qt.AlignHCenter
                    text: "TabBar is a bar with icons or text which allows the user"
                          + "to switch between different subtasks, views, or modes."
                }
            }
        }
    }

    header: TabBar {
        id: tabBar
        TabButton {
            text: qsTr("Account")
        }
        TabButton {
            text: qsTr("Test_2")
        }
        TabButton {
            text: qsTr("Test_3")
        }
    }
}
