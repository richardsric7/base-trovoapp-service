// ignore_for_file: must_be_immutable

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../fonts.dart';
import '../notifire_clor.dart';

class Button extends StatefulWidget {
  final String? buttontext;
  final Color? colorbutton;
  final Color? buttontextcolor;
  final double? width;
  final double? height;
  final void Function()? onTap;

  const Button(
    this.buttontext,
    this.colorbutton,
    this.buttontextcolor, {
    Key? key,
    this.onTap,
    this.height,
    this.width,
  }) : super(key: key);

  @override
  State<Button> createState() => _ButtonState();
}

class _ButtonState extends State<Button> {
  get borderRadius => BorderRadius.circular(15);

  late ColorNotifier notifier;

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return Container(
      decoration: BoxDecoration(
        borderRadius: borderRadius,
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: <Widget>[
          LayoutBuilder(builder: (context, constraints) {
            return Container(
              height: widget.height ?? height / 15,
              width: widget.width ?? width / 1.1,
              child: ElevatedButton(
                onPressed: widget.onTap,
                style: ButtonStyle(
                  backgroundColor:
                      MaterialStateProperty.all<Color>(widget.colorbutton!),
                  shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                    const RoundedRectangleBorder(
                      borderRadius: BorderRadius.all(
                        Radius.circular(15),
                      ),
                    ),
                  ),
                ),
                child: Center(
                  child: Text(
                    widget.buttontext!,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                        fontFamily: fontbody,
                        fontSize: 15,
                        color: widget.buttontextcolor),
                  ),
                ),
              ),
            );
          }),
        ],
      ),
    );
  }
}

class ButtonWithIcon extends StatefulWidget {
  final String? buttontext;
  final Color? colorbutton;
  final Color? buttontextcolor;
  final String imageUrl;
  final double? width;
  final double? height;
  final void Function()? onTap;

  const ButtonWithIcon(
    this.buttontext,
    this.colorbutton,
    this.buttontextcolor,
    this.imageUrl, {
    Key? key,
    this.onTap,
    this.height,
    this.width,
  }) : super(key: key);

  @override
  State<ButtonWithIcon> createState() => _ButtonWithIconState();
}

class _ButtonWithIconState extends State<ButtonWithIcon> {
  get borderRadius => BorderRadius.circular(15);

  late ColorNotifier notifier;

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return Container(
      decoration: BoxDecoration(
        borderRadius: borderRadius,
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: <Widget>[
          LayoutBuilder(builder: (context, constraints) {
            return Container(
              height: widget.height ?? height / 15,
              width: widget.width ?? width / 1.1,
              child: ElevatedButton(
                onPressed: widget.onTap,
                style: ButtonStyle(
                  backgroundColor:
                      MaterialStateProperty.all<Color>(widget.colorbutton!),
                  shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                    const RoundedRectangleBorder(
                      borderRadius: BorderRadius.all(
                        Radius.circular(15),
                      ),
                    ),
                  ),
                ),
                child: Row(
                  children: [
                    Image.asset(
                      widget.imageUrl,
                      height: 50,
                      width: 50,
                    ),
                    SizedBox(
                      width: 5,
                    ),
                    Text(
                      widget.buttontext!,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                          fontFamily: fontbody,
                          fontSize: 15,
                          color: widget.buttontextcolor),
                    ),
                  ],
                ),
              ),
            );
          }),
        ],
      ),
    );
  }
}

class HalfButtonWithIcon extends StatefulWidget {
  final String? buttontext;
  final Color? colorbutton;
  final Color? buttontextcolor;
  final String imageUrl;
  final double? width;
  final double? height;
  final void Function()? onTap;

  const HalfButtonWithIcon(
    this.buttontext,
    this.colorbutton,
    this.buttontextcolor,
    this.imageUrl, {
    Key? key,
    this.onTap,
    this.height,
    this.width,
  }) : super(key: key);

  @override
  State<HalfButtonWithIcon> createState() => _HalfButtonWithIconState();
}

class _HalfButtonWithIconState extends State<HalfButtonWithIcon> {
  get borderRadius => BorderRadius.circular(15);

  late ColorNotifier notifier;

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return Container(
      decoration: BoxDecoration(
        borderRadius: borderRadius,
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.start,
        children: <Widget>[
          LayoutBuilder(builder: (context, constraints) {
            return Container(
              height: widget.height ?? height / 15,
              width: widget.width ?? width / 1.1,
              child: ElevatedButton(
                onPressed: widget.onTap,
                style: ButtonStyle(
                  backgroundColor:
                      MaterialStateProperty.all<Color>(widget.colorbutton!),
                  shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                    const RoundedRectangleBorder(
                      borderRadius: BorderRadius.all(
                        Radius.circular(15),
                      ),
                    ),
                  ),
                ),
                child: Row(
                  children: [
                    Image.asset(
                      widget.imageUrl,
                      height: 40,
                      width: 40,
                    ),
                    Container(
                      constraints: BoxConstraints(maxWidth: 100),
                      child: Text(
                        widget.buttontext!,
                        textAlign: TextAlign.center,
                        style: TextStyle(
                            fontFamily: fontbody,
                            fontSize: 12,
                            color: widget.buttontextcolor),
                      ),
                    ),
                  ],
                ),
              ),
            );
          }),
        ],
      ),
    );
  }
}

class ButtonOutlined extends StatefulWidget {
  final String? buttontext;
  final Color? colorbutton;
  final Color? buttontextcolor;
  final Color? borderColor;
  final double? width;
  final double? height;
  final void Function()? onTap;

  const ButtonOutlined(this.buttontext, this.colorbutton, this.buttontextcolor,
      {Key? key, this.onTap, this.width, this.height, this.borderColor})
      : super(key: key);

  @override
  State<ButtonOutlined> createState() => _ButtonOutlinedState();
}

class _ButtonOutlinedState extends State<ButtonOutlined> {
  get borderRadius => BorderRadius.circular(15);

  late ColorNotifier notifier;

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return Container(
      decoration: BoxDecoration(
        borderRadius: borderRadius,
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          LayoutBuilder(builder: (context, constraints) {
            return Container(
              height: widget.height ?? height / 15,
              width: widget.width ?? width / 1.1,
              child: TextButton(
                onPressed: widget.onTap,
                style: ButtonStyle(
                  overlayColor:
                      MaterialStateProperty.all<Color>(notifier.getsplashgrey),
                  elevation: MaterialStateProperty.all<double>(0),
                  backgroundColor:
                      MaterialStateProperty.all<Color>(widget.colorbutton!),
                  side: MaterialStateProperty.all(
                    BorderSide(
                        color: widget.borderColor ?? notifier.getgrey,
                        width: 1,
                        style: BorderStyle.solid),
                  ),
                  shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                    const RoundedRectangleBorder(
                      borderRadius: BorderRadius.all(
                        Radius.circular(10),
                      ),
                    ),
                  ),
                ),
                child: Center(
                  child: Text(
                    widget.buttontext!,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                        fontFamily: fontbody,
                        fontSize: 15,
                        color: widget.buttontextcolor),
                  ),
                ),
              ),
            );
          }),
        ],
      ),
    );
  }
}

class SmallButtonOutlined extends StatefulWidget {
  final String? buttontext;
  final Color? colorbutton;
  final Color? buttontextcolor;
  final void Function()? onTap;

  const SmallButtonOutlined(
      this.buttontext, this.colorbutton, this.buttontextcolor,
      {Key? key, this.onTap})
      : super(key: key);

  @override
  State<SmallButtonOutlined> createState() => _SmallButtonOutlinedState();
}

class _SmallButtonOutlinedState extends State<SmallButtonOutlined> {
  get borderRadius => BorderRadius.circular(15);

  late ColorNotifier notifier;

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ElevatedButton(
      onPressed: widget.onTap,
      style: ButtonStyle(
        overlayColor: MaterialStateProperty.all<Color>(notifier.getsplashgrey),
        backgroundColor: MaterialStateProperty.all<Color>(widget.colorbutton!),
        side: MaterialStateProperty.all(
          BorderSide(
              color: notifier.getbluewhitecolor,
              width: 1,
              style: BorderStyle.solid),
        ),
        shape: MaterialStateProperty.all<RoundedRectangleBorder>(
          const RoundedRectangleBorder(
            borderRadius: BorderRadius.all(
              Radius.circular(10),
            ),
          ),
        ),
      ),
      child: Text(
        widget.buttontext!,
        style: TextStyle(
          fontFamily: fontsemibold,
          color: widget.buttontextcolor,
        ),
      ),
    );
  }
}

class SmallButton extends StatefulWidget {
  final String? buttontext;
  final Color? colorbutton;
  final Color? buttontextcolor;
  final void Function()? onTap;

  const SmallButton(
    this.buttontext,
    this.colorbutton,
    this.buttontextcolor, {
    Key? key,
    this.onTap,
  }) : super(key: key);

  @override
  State<SmallButton> createState() => _SmallButtonState();
}

class _SmallButtonState extends State<SmallButton> {
  get borderRadius => BorderRadius.circular(15);

  late ColorNotifier notifier;

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ElevatedButton(
      onPressed: widget.onTap,
      style: ButtonStyle(
        backgroundColor: MaterialStateProperty.all<Color>(widget.colorbutton!),
        shape: MaterialStateProperty.all<RoundedRectangleBorder>(
          const RoundedRectangleBorder(
            borderRadius: BorderRadius.all(
              Radius.circular(10),
            ),
          ),
        ),
      ),
      child: Text(
        widget.buttontext!,
        style: TextStyle(
          fontFamily: fontbody,
          fontSize: 15,
          color: widget.buttontextcolor,
        ),
      ),
    );
  }
}
