package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"modelconv/converter"
	"modelconv/mqo"
)

func defaultOutputFile(input string) string {
	ext := strings.ToLower(filepath.Ext(input))
	base := input[0 : len(input)-len(ext)]
	if ext == ".glb" {
		return base + ".vrm"
	} else if ext == ".mqo" {
		return base + ".glb"
	} else if ext == ".pmx" || ext == ".pmd" {
		return base + ".mqo"
	}
	return input + ".mqo"
}

func saveDocument(doc *mqo.Document, output, srcDir, vrmConf string, forceUnlit bool) error {
	ext := strings.ToLower(filepath.Ext(output))
	if ext == ".glb" {
		options := &converter.MQOToGLTFOption{ForceUnlit: forceUnlit}
		return saveAsGlb(doc, output, srcDir, options)
	} else if ext == ".vrm" {
		options := &converter.MQOToGLTFOption{ForceUnlit: forceUnlit}
		gltfdoc, err := converter.NewMQOToGLTFConverter(options).Convert(doc, srcDir)
		if err != nil {
			return err
		}
		return saveAsVRM(gltfdoc, output, vrmConf)
	} else if ext == ".mqo" {
		w, err := os.Create(output)
		if err != nil {
			return err
		}
		defer w.Close()
		return mqo.WriteMQO(doc, w, output)
	} else if ext == ".pmx" {
		return saveAsPmx(doc, output)
	}
	return fmt.Errorf("Unsuppored output type: %v", ext)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s input.pmx [output.mqo]\n", os.Args[0])
		flag.PrintDefaults()
	}
	rot180 := flag.Bool("rot180", false, "rotate 180 degrees around Y (.mqo)")
	autoTpose := flag.String("autotpose", "", "Arm bone names(.mqo)")
	forceUnlit := flag.Bool("gltfunlit", false, "unlit all materials")
	scale := flag.Float64("scale", 0, "0:auto")
	vrmconf := flag.String("vrmconfig", "", "config file for VRM")
	hides := flag.String("hide", "", "hide objects")
	hidemats := flag.String("hidemat", "", "hide materials")
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}
	input := flag.Arg(0)
	output := flag.Arg(1)
	if output == "" {
		output = defaultOutputFile(input)
	}
	confFile := *vrmconf
	if confFile == "" {
		confFile = input[0:len(input)-len(filepath.Ext(input))] + ".vrmconfig.json"
		if _, err := os.Stat(confFile); err != nil {
			confFile = ""
		}
	}

	inputExt := strings.ToLower(filepath.Ext(input))
	if inputExt == ".glb" {
		err := glb2vrm(input, output, confFile)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *scale == 0 {
		if inputExt == ".pmx" || inputExt == ".pmd" {
			*scale = 80
		} else {
			*scale = 1
		}
	}

	// mmd to vrm
	outputExt := strings.ToLower(filepath.Ext(output))
	if (inputExt == ".pmx" || inputExt == ".pmd") && outputExt == ".vrm" {
		*rot180 = true
		if confFile == "" {
			execPath, _ := os.Executable()
			confFile = filepath.Join(filepath.Dir(execPath), "vrmconfig_presets/mmd_ja.json")
		}
		if *autoTpose == "" {
			*autoTpose = "右腕,左腕"
		}
	}

	doc, err := loadDocument(input)
	if err != nil {
		log.Fatal(err)
	}

	// transform
	if *rot180 {
		doc.Transform(func(v *mqo.Vector3) {
			v.X *= -1
			v.Z *= -1
		})
	}
	if *scale != 1.0 {
		s := float32(*scale)
		doc.Transform(func(v *mqo.Vector3) {
			v.X *= s
			v.Y *= s
			v.Z *= s
		})
	}

	if *autoTpose != "" {
		for _, boneName := range strings.Split(*autoTpose, ",") {
			for _, b := range mqo.GetBonePlugin(doc).Bones() {
				if b.Name == boneName {
					doc.BoneAdjustX(b)
				}
			}
		}
	}

	if *hides != "" {
		objectByName := map[string]int{}
		for idx, obj := range doc.Objects {
			objectByName[obj.Name] = idx
		}
		for _, n := range strings.Split(*hides, ",") {
			if idx, ok := objectByName[n]; ok {
				d := doc.Objects[idx].Depth
				doc.Objects[idx].Visible = false
				idx++
				for idx < len(doc.Objects) && doc.Objects[idx].Depth > d {
					doc.Objects[idx].Visible = false
					idx++
				}
			}
		}
	}

	if *hidemats != "" {
		materialsByName := map[string]int{}
		for idx, mat := range doc.Materials {
			materialsByName[mat.Name] = idx
		}
		for _, n := range strings.Split(*hidemats, ",") {
			if idx, ok := materialsByName[n]; ok {
				doc.Materials[idx].Name = "$IGNORE"
			}
		}
	}

	log.Print("out: ", output)
	if err = saveDocument(doc, output, filepath.Dir(input), confFile, *forceUnlit); err != nil {
		log.Fatal(err)
	}
}
